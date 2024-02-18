package src

import (
	"encoding/json"
	"fmt"
	"marketplace_server/config"
	"marketplace_server/internal/servers"

	model_backpack "marketplace_server/internal/backpack/model"
	model_transaction "marketplace_server/internal/bill/model"
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/common/rabbitmqx"
	"marketplace_server/internal/common/utils"
	"marketplace_server/internal/product/Infrastructure_layer"
	model_product "marketplace_server/internal/product/model"

	// /Users/liuming/Documents/Work/DDD/barketplace_server/internal/bill/model/transactionl_entity.go
	"marketplace_server/internal/user/model"
	"runtime/debug"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

const (
	ExchangeType           = "direct"
	TransactionExchange    = "transaction_exchange"        // 通知交换机
	BindKeyPurchaseProduct = "notify_purchase_product_key" // 通用邮件绑定key
)

// 交易引擎
type TransactionEgine struct {
	DataLock sync.RWMutex
	cfg      *config.SugaredConfig

	Repos *servers.RepositoriesManager // 持久層管理

	PurchaseProductList []*model.ProductTransactionParams // 購買等候清單 會選slice 是因為 元素越小優先越高, 可重複快速搜尋
	SellProductList     []*model.ProductTransactionParams // 販賣等候清單
	marketPriceMap      map[string]string                 // 市場最新價格 key=商品名稱 value={"product_count":1000,"currency":"TWD","amount":"10"}
	SysRate             decimal.Decimal                   // 系統抽成
	Consumer            *rabbitmqx.Consumer               // mq
}

// 建立交易引擎
func NewTransactionEgine(cfg *config.SugaredConfig) *TransactionEgine {

	// 建立 db連線 和 redis連線
	repos := servers.NewRepositories(cfg)
	repos.Automigrate()

	// 綁定交易搓合物件
	transactionEgine := &TransactionEgine{
		cfg:            cfg,
		Repos:          repos,                     // 持久層
		marketPriceMap: make(map[string]string),   // 市場價格
		SysRate:        decimal.NewFromFloat(1.0), // 系統抽成, 目前沒抽
	}

	logs.Debugf("RFC3339 start time:%v", time.Now().Format(time.RFC3339))

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
		return nil
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

	// 資料鎖
	t.DataLock.Lock()
	defer t.DataLock.Unlock()

	// 沒有資料就不用搓合
	if len(t.PurchaseProductList) == 0 {
		return
	}
	if len(t.SellProductList) == 0 {
		return
	}

	// 撈取市場最新價格 (取得redis緩存)
	dataMap, err := t.Repos.ProductRepo.RedisGetMarketPrice(Infrastructure_layer.Redis_MarketPrice)
	if err != nil {
		return
	}
	t.marketPriceMap = dataMap
	logs.Debugf("marketPriceMap:%+v", t.marketPriceMap)

	// 搜尋優先配對搓合的 購買清單
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

		// 取得買方的價格
		purchaseAmount := purchaseData.GetPrice(marketPriceDetail.Amount)

		logs.Debugf("i:%d, amount(買的價格):%s, marketPriceDetail(市場價格):%+v",
			i, purchaseAmount.String(), marketPriceDetail)

		// 搜尋優先配對搓合的 販賣清單
		for j, sellData := range t.SellProductList {

			// 比對產品是否相同
			if purchaseData.ProductName != sellData.ProductName {
				continue
			}
			// 比對 相同用戶 不給予搓則
			if purchaseData.UserID == sellData.UserID {
				continue
			}

			isGet := false

			// 取得賣方想要的價格
			sellAmount := sellData.GetPrice(marketPriceDetail.Amount)
			// var sellAmount decimal.Decimal
			// switch model.TransferType(sellData.TransferType) {
			// case model.LimitPrice: // 限價
			// 	//data.Amount(現價的價格) >= price(市場價格) 才能買到
			// 	// ret := purchaseData.Amount.GreaterThanOrEqual(marketPriceDetail.Amount)
			// 	// if ret {
			// 	// 	// 配對成功
			// 	// 	logs.Debugf("#### 配對成功")
			// 	// }

			// 	// 取得賣方的現價 價格
			// 	sellAmount = sellData.Amount
			// case model.MarketPrice: // 市價

			// 	// 使用市場價格當賣方價格
			// 	sellAmount = marketPriceDetail.Amount

			// 	// 如果 賣方價格
			// 	// price := t.marketPriceMap[data.ProductName]
			// 	// data.Amount >= price
			// default:
			// 	logs.Warnf("錯誤的 transferType data:%+v", purchaseData)
			// 	continue
			// }
			logs.Debugf("j:%d, amount(賣的價格):%s, marketPriceDetail(市場價格):%+v",
				j, sellAmount.String(), marketPriceDetail)

			// 如果 買方價格 >= 賣方
			isGet = purchaseAmount.GreaterThanOrEqual(sellAmount)
			logs.Debugf("配對開始 isGet:%v 買:%v >= 賣:%v",
				isGet, purchaseAmount.String(), sellAmount.String())

			if isGet {

				// 配對成功
				logs.Debugf(" #### 配對成功 買:%v >= 賣:%v",
					purchaseData.Amount.String(), sellAmount.String())

				// 寫進db (使用 transaction(事務) 失敗就Rollback)

				// 寫入買方背包內
				backpackObj := &model_backpack.Backpack{
					UserID:       purchaseData.UserID, // 買方用戶ID
					ProductName:  purchaseData.ProductName,
					ProductCount: purchaseData.OperateCount,
					CreatedAt:    time.Now(), // 創建時間
					UodateAt:     time.Now(), // 更新時間
				}
				t.Repos.BackpackRepo.Save(backpackObj)

				// 扣除賣方商品的數量

				// 使用賣方的價格當作成交價, 更新賣家交易單
				sellTransaction, err := t.Repos.TransactionRepo.GetTransactionInfo(sellData.TransactionID)
				if err != nil {
					logs.Warnf("getTransactionInfo transactionID:%v, err:%v", sellData.TransactionID, err)
					continue
				}
				sellTransaction.Amount = sellAmount                                        // 更新交易價格
				sellTransaction.UodateAt = time.Now()                                      // 更新交易完成時間
				sellTransaction.ToUserID = purchaseData.UserID                             // 買家的id
				sellTransaction.Status = int8(model_transaction.Transaction_Status_Finish) // 交易完成狀態
				t.Repos.TransactionRepo.Save(sellTransaction)

				// 使用賣方的價格當作成交價, 更新買家交易單
				purchaseTransaction, err := t.Repos.TransactionRepo.GetTransactionInfo(purchaseData.TransactionID)
				if err != nil {
					logs.Warnf("getTransactionInfo transactionID:%v, err:%v", purchaseData.TransactionID, err)
					continue
				}
				purchaseTransaction.Amount = sellAmount                                        // 更新交易價格
				purchaseTransaction.UodateAt = time.Now()                                      // 更新交易完成時間
				purchaseTransaction.ToUserID = sellData.UserID                                 // 賣家的id
				purchaseTransaction.Status = int8(model_transaction.Transaction_Status_Finish) // 交易完成狀態
				t.Repos.TransactionRepo.Save(purchaseTransaction)

				// 更新買家用戶金額 = 買家目前金額 - 賣家金額
				purchaseUser, err := t.Repos.UserRepo.GetUserInfo(purchaseData.UserID)
				if err != nil {
					logs.Warnf("getTransactionInfo transactionID:%v, err:%v", purchaseData.UserID, err)
					continue
				}
				purchaseUser.Amount.Sub(sellAmount)
				_, err = t.Repos.UserRepo.Save(purchaseUser)
				if err != nil {
					logs.Warnf("getTransactionInfo transactionID:%v, err:%v", purchaseData.UserID, err)
					continue
				}

				// 更新賣家用戶的金額 = 賣家用戶的金額 + (販賣 * 系統抽成)
				sellUser, err := t.Repos.UserRepo.GetUserInfo(sellData.UserID)
				if err != nil {
					logs.Warnf("getTransactionInfo transactionID:%v, err:%v", sellData.UserID, err)
					continue
				}
				sellUser.Amount.Add(sellAmount.Mul(t.SysRate))
				_, err = t.Repos.UserRepo.Save(sellUser)
				if err != nil {
					logs.Warnf("getTransactionInfo transactionID:%v, err:%v", sellData.UserID, err)
					continue
				}

				// 更新回redis, 市場最新價格 例如 t.marketPriceMap["BTC"] = 賣方價格 元成交
				marketPriceDetail.Amount = sellAmount
				updateMarketPriceJson, _ := json.Marshal(marketPriceDetail.Amount)
				t.marketPriceMap[sellData.ProductName] = string(updateMarketPriceJson)
				err = t.Repos.ProductRepo.RedisSetMarketPrice(Infrastructure_layer.Redis_MarketPrice, t.marketPriceMap)
				if err != nil {
					return
				}

				// 寄送mq 給 marketplace_server

				// 刪除 配對搓合的購買清單
				utils.SliceHelper(&t.PurchaseProductList).Remove(i)
				utils.SliceHelper(&t.SellProductList).Remove(j)
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
		logs.Errorf("dispatch fail productTransactionNotify:%+v, err:%v",
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
