package src

import (
	"encoding/json"
	"marketplace_server/config"
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/common/rabbitmqx"
	"marketplace_server/internal/user/model"
	"time"
)

const (
	ExchangeType           = "direct"
	TransactionExchange    = "transaction_exchange"        // 通知交换机
	BindKeyPurchaseProduct = "notify_purchase_product_key" // 通用邮件绑定key
)

// 交易引擎
type TransactionEgine struct {
	cfg *config.SugaredConfig
}

// 建立交易引擎
func NewTransactionEgine(cfg *config.SugaredConfig) *TransactionEgine {

	transactionEgine := &TransactionEgine{
		cfg: cfg,
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

// 收到交易通知
func (t *TransactionEgine) NotifyTransaction(message []byte) error {

	//logs.Debugf("msg:%s", string(message[:]))

	productTransactionNotify := &model.ProductTransactionNotify{}
	err := json.Unmarshal(message, productTransactionNotify)
	if err != nil {
		logs.Errorf("unmarshal err, err:%v, message:%v", err, string(message[:]))
		return nil
	}

	logs.Debugf("productTransactionNotify:%+v", productTransactionNotify)

	// cmd 分配
	switch productTransactionNotify.Cmd {
	case model.Notify_Cmd_Purchase:
	case model.Notify_Cmd_Sell:
	case model.Notify_Cmd_Cancel:
	default:
		logs.Warnf("unkonw cmd:%v", productTransactionNotify.Cmd)
	}

	return nil
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

	logs.Debugf("啟動搓合監聽goroutine")
	for {

		time.Sleep(time.Second * 30) // 30秒搓合一次
	}
}
