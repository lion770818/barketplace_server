package src

import (
	"marketplace_server/config"
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/common/rabbitmqx"
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

	msg := string(message[:])
	logs.Debugf("msg:%s", msg)

	// mailMessage := &notify.ReqMailGeneralParam{}
	// err := json.Unmarshal(_message, mailMessage)
	// if err != nil {
	// 	logs.Errorf("unmarshal err, err:%v, message:%v", err, _message)
	// 	return nil
	// }

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
