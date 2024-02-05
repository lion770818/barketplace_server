package src

import (
	"encoding/json"
	"fmt"
	"marketplace_server/config"
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/common/rabbitmqx"
	"marketplace_server/internal/common/utils"
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
	DataLock sync.RWMutex
	cfg      *config.SugaredConfig

	PurchaseList   []*model.ProductPurchaseParams // 會選slice 是因為 元素越小優先越高, 可重複快速搜尋
	marketPriceMap map[string]string              // 市場最新價格 key=商品名稱 value={"product_count":1000,"currency":"TWD","amount":"10"}

}

// 建立交易引擎
func NewTransactionEgine(cfg *config.SugaredConfig) *TransactionEgine {

	transactionEgine := &TransactionEgine{
		cfg:            cfg,
		marketPriceMap: make(map[string]string),
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
func (t *TransactionEgine) consumeNotifyTransaction(_host, _port, _user, _password string, _connectionNum, _channelNum int, tag string) {

	uri := "amqp://" + _user + ":" + _password + "@" + _host + ":" + _port + "/"

	consumer := rabbitmqx.NewConsumer(uri, ExchangeType,
		TransactionExchange, BindKeyPurchaseProduct,
		BindKeyPurchaseProduct, tag, false,
		true, t.NotifyTransaction)
	if err := consumer.Start(); err != nil {
		logs.Errorf("RabbitInit error,err = " + err.Error())
	}
}

func (t *TransactionEgine) Run() {

	defer func() {

		if err := recover(); err != nil {
			logs.Warnf("引發例外 err:%+v", err, string(debug.Stack()))
		}
	}()

	logs.Debugf("啟動搓合監聽goroutine")
	for {

		time.Sleep(time.Second * 30) // 30秒搓合一次
	}
}

// 排程任務
func (t *TransactionEgine) Cron() {

	// 	資料鎖
	t.DataLock.Lock()
	defer t.DataLock.Unlock()

	// 撈取市場最新價格
	// t.marketPriceMap = redis 來的市場價格資料

	// 搜尋優先配對搓合的購買清單
	for i, data := range t.PurchaseList {

		isGet := false
		switch model.TransferType(data.TransferType) {

		case model.LimitPrice: // 限價
		case model.MarketPrice: // 市價
		default:
			logs.Warnf("錯誤的 transferType data:%+v", data)
			continue
		}

		// todo 假設找到想配對的清單
		if isGet {

			// 寫進db (使用 transaction(事務) 失敗就Rollback)

			// 更新回redis, 市場最新價格 例如 t.marketPriceMap["BTC"] = 20 元成交

			// 寄送mq 給 marketplace_server

			// 刪除 配對搓合的購買清單
			utils.SliceHelper(&data).Remove(i)
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
		err = t.Purchase(productTransactionNotify)
	case model.Notify_Cmd_Sell:
	case model.Notify_Cmd_Cancel:
	default:
		logs.Warnf("unkonw cmd:%v", productTransactionNotify.Cmd)
	}

	return
}

// 交易 買
func (t *TransactionEgine) Purchase(productTransactionNotify *model.ProductTransactionNotify) error {

	// 解析封包
	// var productPurchaseParams model.ProductPurchaseParams
	// err := json.Unmarshal(productTransactionNotify.Data.([]byte), &productPurchaseParams)
	// if err != nil {
	// 	return err
	// }

	// // 寫入 買結構
	// t.PurchaseList = append(t.PurchaseList, &productPurchaseParams)

	// logs.Debugf("等待購買清單:%d, 新進資料:%+v", len(t.PurchaseList), productPurchaseParams)
	return nil
}

// 交易 賣
func (t *TransactionEgine) Sell(productTransactionNotify *model.ProductTransactionNotify) error {

	// 解析封包

	return nil
}

// 取消交易
func (t *TransactionEgine) Cancel(productTransactionNotify *model.ProductTransactionNotify) error {

	// 解析封包

	// 搜尋要取消的清單
	for i, data := range t.PurchaseList {

		// todo 假設找到想取消的清單
		if true {
			utils.SliceHelper(&data).Remove(i)
		}
	}

	return nil
}
