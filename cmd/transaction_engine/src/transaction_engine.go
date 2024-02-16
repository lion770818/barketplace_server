package src

import (
	"encoding/json"
	"fmt"
	"marketplace_server/config"
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/common/mysql"
	"marketplace_server/internal/common/rabbitmqx"
	"marketplace_server/internal/common/redis"
	"marketplace_server/internal/common/utils"
	"marketplace_server/internal/product/Infrastructure_layer"
	Infrastructure_product "marketplace_server/internal/product/Infrastructure_layer"
	model_product "marketplace_server/internal/product/model"
	"marketplace_server/internal/user/model"
	"runtime/debug"
	"sync"
	"time"
)

const (
	ExchangeType           = "direct"
	TransactionExchange    = "transaction_exchange"        // 通知交换机
	BindKeyPurchaseProduct = "notify_purchase_product_key" // 通用邮件绑定key
)

// 交易引擎
type TransactionEgine struct {
	DataLock    sync.RWMutex
	cfg         *config.SugaredConfig
	ProductRepo Infrastructure_layer.ProductRepo // 產品持久層

	PurchaseProductList []*model.ProductTransactionParams // 購買等候清單 會選slice 是因為 元素越小優先越高, 可重複快速搜尋
	SellProductList     []*model.ProductTransactionParams // 販賣等候清單
	marketPriceMap      map[string]string                 // 市場最新價格 key=商品名稱 value={"product_count":1000,"currency":"TWD","amount":"10"}

	Consumer *rabbitmqx.Consumer // mq
}

// 建立交易引擎
func NewTransactionEgine(cfg *config.SugaredConfig) *TransactionEgine {

	// 持久化类型的 repo
	mysqlCfg := mysql.Config{
		LogMode:  cfg.Mysql.LogMode,
		Driver:   cfg.Mysql.Driver,
		Host:     cfg.Mysql.Host,
		Port:     cfg.Mysql.Port,
		Database: cfg.Mysql.Database,
		User:     cfg.Mysql.User,
		Password: cfg.Mysql.Password,
	}
	db := mysql.NewDB(mysqlCfg)

	// 建立redis連線
	redisCfg := &redis.RedisParameter{
		Network:      "tcp",
		Address:      fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           0,
		DialTimeout:  time.Second * time.Duration(10),
		ReadTimeout:  time.Second * time.Duration(10),
		WriteTimeout: time.Second * time.Duration(10),
		PoolSize:     10,
	}
	redisClient, err := redis.NewRedis(redisCfg)
	if err != nil {
		logs.Errorf("newRedis error=%v", err)
		return nil
	}

	protuctRepo := Infrastructure_product.NewProductRepoManager(db, redisClient.GetClient())

	transactionEgine := &TransactionEgine{
		cfg:            cfg,
		marketPriceMap: make(map[string]string),
		ProductRepo:    protuctRepo,
	}

	// 監聽 rabbit mq
	transactionEgine.consumeNotifyTransaction(cfg.RabbitMq.Host,
		cfg.RabbitMq.Port,
		cfg.RabbitMq.User,
		cfg.RabbitMq.Password,
		cfg.RabbitMq.ConnectNum,
		cfg.RabbitMq.ChannelNum,
		"test")
	return transactionEgine
}

// 建立 交易通知 的 消費端
func (t *TransactionEgine) consumeNotifyTransaction(_host, _port, _user, _password string, _connectionNum, _channelNum int, tag string) (consumer *rabbitmqx.Consumer) {

	uri := "amqp://" + _user + ":" + _password + "@" + _host + ":" + _port + "/"

	consumer = rabbitmqx.NewConsumer(uri, ExchangeType,
		TransactionExchange, BindKeyPurchaseProduct,
		BindKeyPurchaseProduct, tag, false,
		true, t.NotifyTransaction)
	if err := consumer.Start(); err != nil {
		logs.Errorf("RabbitInit error,err = " + err.Error())
	}

	return consumer
}

func (t *TransactionEgine) Run() {

	defer func() {

		if err := recover(); err != nil {
			logs.Warnf("引發例外 err:%+v", err, string(debug.Stack()))
		}
	}()

	logs.Debugf("啟動搓合監聽goroutine")
	for {

		t.Cron()                     // 搓合
		time.Sleep(time.Second * 30) // 30秒搓合一次
	}
}

// 排程任務
func (t *TransactionEgine) Cron() {

	// 	資料鎖
	t.DataLock.Lock()
	defer t.DataLock.Unlock()

	// 沒有資料就不用搓合
	if len(t.PurchaseProductList) == 0 {
		return
	}
	if len(t.SellProductList) == 0 {
		return
	}

	//  撈取市場最新價格 (取得redis緩存)
	dataMap, err := t.ProductRepo.RedisGetMarketPrice(Infrastructure_layer.Redis_MarketPrice)
	if err != nil {
		return
	}
	t.marketPriceMap = dataMap
	logs.Debugf("marketPriceMap:%+v", t.marketPriceMap)

	// 搜尋優先配對搓合的購買清單
	for i, purchaseData := range t.PurchaseProductList {

		// 取得要配對的商品的市場價格
		productName := purchaseData.ProductName
		marketPriceJson, ok := t.marketPriceMap[productName]
		if !ok {
			logs.Warnf("快取不存在的產品 productName:%v, marketPriceMap:%+v",
				productName, t.marketPriceMap)
			continue
		}

		// 取得市場價格物件
		marketPriceDetail, err := model_product.NewMarketPriceRedis(marketPriceJson)
		if err != nil {
			logs.Warnf("marketPriceJson:%v, err:%v", marketPriceJson, err)
			continue
		}

		logs.Debugf("i:%d, amount(買的價格):%s, marketPriceDetail(市場價格):%+v",
			i, purchaseData.Amount.String(), marketPriceDetail)

		// 搜尋優先配對搓合的販賣清單
		for j, sellData := range t.SellProductList {

			logs.Debugf("j:%d, amount(賣的價格):%s, marketPriceDetail(市場價格):%+v",
				j, sellData.Amount.String(), marketPriceDetail)
			isGet := false
			switch model.TransferType(purchaseData.TransferType) {

			case model.LimitPrice: // 限價
				//data.Amount(現價的價格) >= price(市場價格) 才能買到
				ret := purchaseData.Amount.GreaterThanOrEqual(marketPriceDetail.Amount)
				if ret {
					// 配對成功
					logs.Debugf("#### 配對成功")
				}
			case model.MarketPrice: // 市價

				// 如果 賣方價格
				// price := t.marketPriceMap[data.ProductName]
				// data.Amount >= price
			default:
				logs.Warnf("錯誤的 transferType data:%+v", purchaseData)
				continue
			}

			// todo 假設找到想配對的清單
			if isGet {

				// 寫進db (使用 transaction(事務) 失敗就Rollback)

				// 更新回redis, 市場最新價格 例如 t.marketPriceMap["BTC"] = 20 元成交

				// 寄送mq 給 marketplace_server

				// 刪除 配對搓合的購買清單
				utils.SliceHelper(&purchaseData).Remove(i)
			}
		}

	}
}

// 收到交易通知
func (t *TransactionEgine) NotifyTransaction(message []byte) error {

	logs.Debugf("msg:%s", string(message[:]))

	// 去除斜線 轉成 map
	dataMap := map[string]interface{}{}
	err := json.Unmarshal(message, &dataMap)
	if err != nil {
		logs.Errorf("unmarshal err, err:%v, message:%v", err, string(message[:]))
		return nil
	}

	dataMapTmp, err := json.Marshal(dataMap)
	if err != nil {
		logs.Errorf("marshal err, err:%v, message:%v", err, string(message[:]))
		return nil
	}
	// map to obj
	productTransactionNotify := &model.ProductTransactionNotify{}
	err = json.Unmarshal(dataMapTmp, productTransactionNotify)
	if err != nil {
		logs.Errorf("unmarshal err, err:%v, message:%v", err, string(message[:]))
		return nil
	}

	logs.Debugf("productTransactionNotify:%+v", productTransactionNotify)

	// 封包分派
	err = t.Dispatch(productTransactionNotify)
	if err != nil {
		logs.Debugf("dispatch fail productTransactionNotify:%+v, err:%v",
			productTransactionNotify, err)
	}

	return nil
}

// 封包分派
func (t *TransactionEgine) Dispatch(productTransactionNotify *model.ProductTransactionNotify) (err error) {

	// 	資料鎖
	t.DataLock.Lock()
	defer t.DataLock.Unlock()

	if productTransactionNotify == nil {
		err = fmt.Errorf("productTransactionNotify == nil")
		return
	}

	// cmd 分配
	switch productTransactionNotify.Cmd {
	case model.Notify_Cmd_Purchase:
		err = t.PurchaseProduct(productTransactionNotify) // 儲存到購買清單
	case model.Notify_Cmd_Sell:
		err = t.SellProduct(productTransactionNotify) // 儲存到販賣清單
	case model.Notify_Cmd_Cancel:
	default:
		logs.Warnf("unkonw cmd:%v", productTransactionNotify.Cmd)
	}

	return
}

// 交易 買
func (t *TransactionEgine) PurchaseProduct(productTransactionNotify *model.ProductTransactionNotify) error {

	byteArray, err := json.Marshal(productTransactionNotify.Data)
	if err != nil {
		return err
	}
	// 解析封包
	var productPurchaseParams model.ProductTransactionParams
	err = json.Unmarshal(byteArray, &productPurchaseParams)
	if err != nil {
		return err
	}

	// 模式檢查, 這邊只處理 買
	if model.TransferMode(productPurchaseParams.TransferMode) != model.Purchase {
		return fmt.Errorf("error transaction_mode:%d", productPurchaseParams.TransferMode)
	}

	// 寫入 買結構
	t.PurchaseProductList = append(t.PurchaseProductList, &productPurchaseParams)

	logs.Debugf("等待購買清單:%d, 價格:%s 新進詳細資料:%+v",
		len(t.PurchaseProductList), productPurchaseParams.Amount.String(), productPurchaseParams)
	return nil
}

// 交易 賣
func (t *TransactionEgine) SellProduct(productTransactionNotify *model.ProductTransactionNotify) error {

	// 解析封包 interface to byteArray
	byteArray, err := json.Marshal(productTransactionNotify.Data)
	if err != nil {
		return err
	}
	// 解析封包 byteArray to obj
	var productPurchaseParams model.ProductTransactionParams
	err = json.Unmarshal(byteArray, &productPurchaseParams)
	if err != nil {
		return err
	}

	// 模式檢查, 這邊只處理 賣
	if model.TransferMode(productPurchaseParams.TransferMode) != model.Sell {
		return fmt.Errorf("error transaction_mode:%d", productPurchaseParams.TransferMode)
	}

	// 寫入 賣結構
	t.SellProductList = append(t.SellProductList, &productPurchaseParams)

	logs.Debugf("等待販賣清單:%d, 價格:%s 新進詳細資料:%+v",
		len(t.SellProductList), productPurchaseParams.Amount.String(), productPurchaseParams)
	return nil
}

// 取消交易
func (t *TransactionEgine) Cancel(productTransactionNotify *model.ProductTransactionNotify) error {

	// 解析封包
	byteArray, err := json.Marshal(productTransactionNotify.Data)
	if err != nil {
		return err
	}
	// 解析封包
	var productPurchaseParams model.ProductTransactionParams
	err = json.Unmarshal(byteArray, &productPurchaseParams)
	if err != nil {
		return err
	}

	// 模式檢查, 這邊只處理 賣
	if model.TransferMode(productPurchaseParams.TransferMode) != model.Sell {
		return fmt.Errorf("error transaction_mode:%d", productPurchaseParams.TransferMode)
	}

	// 判斷買或賣
	var waitProductList []*model.ProductTransactionParams
	switch model.TransferMode(productPurchaseParams.TransferMode) {
	case model.Sell:
		waitProductList = t.SellProductList
	case model.Purchase:
		waitProductList = t.PurchaseProductList
	default:
		return fmt.Errorf("error cancel transaction_mode:%d", productPurchaseParams.TransferMode)
	}

	// 搜尋要取消的清單
	for i, data := range waitProductList {

		// todo 假設找到想取消的清單
		if true {
			utils.SliceHelper(&data).Remove(i)
		}
	}

	return nil
}
