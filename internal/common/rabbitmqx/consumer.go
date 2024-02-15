package rabbitmqx

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	quit          chan struct{}
	uri           string
	exchangeType  string
	exchange      string
	queue         string
	routingKey    string
	consumerTag   string
	autoDelete    bool
	durable       bool
	handler       func([]byte) error
}

func NewConsumer(uri, exchangeType, exchange, queue, routingKey, consumerTag string, autoDelete, durable bool, handler func([]byte) error) *Consumer {
	c := &Consumer{
		uri:          uri,
		exchangeType: exchangeType,
		exchange:     exchange,
		queue:        queue,
		routingKey:   routingKey,
		consumerTag:  consumerTag,
		autoDelete:   autoDelete,
		durable:      durable,
		handler:      handler,
		quit:         make(chan struct{}),
	}

	return c
}

func (c *Consumer) Start() error {
	if err := c.Run(); err != nil {
		return err
	}
	go c.ReConnect()

	return nil
}

func (c *Consumer) Stop() {
	close(c.quit)

	if !c.conn.IsClosed() {
		// 关闭 SubMsg message delivery
		if err := c.channel.Cancel(c.consumerTag, true); err != nil {
			log.Println("rabbitmq consumer - channel cancel failed: ", err)
		}

		if err := c.conn.Close(); err != nil {
			log.Println("rabbitmq consumer - connection close failed: ", err)
		}
	}
}

func (c *Consumer) Run() error {
	var err error
	if c.conn, err = amqp.Dial(c.uri); err != nil {
		return err
	}

	if c.channel, err = c.conn.Channel(); err != nil {
		_ = c.conn.Close()
		return err
	}

	log.Printf("exchange:%v, routingKey:%v, queue name:%v", c.exchange, c.routingKey, c.queue)

	if err = c.channel.ExchangeDeclare(
		c.exchange,   // name
		"direct",     // type
		c.durable,    // durable
		c.autoDelete, // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		_ = c.channel.Close()
		_ = c.conn.Close()
		return err
	}

	// declare sms queue
	if _, err = c.channel.QueueDeclare(
		c.queue,      // name
		c.durable,    // durable
		c.autoDelete, // delete when usused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		_ = c.channel.Close()
		_ = c.conn.Close()
		return err
	}

	log.Printf("routingKey:%v, exchange:%v", c.routingKey, c.exchange)
	// bind queue to exchagne by key
	if err = c.channel.QueueBind(
		c.queue,
		c.routingKey,
		c.exchange,
		false,
		nil,
	); err != nil {
		_ = c.channel.Close()
		_ = c.conn.Close()
		return err
	}

	// comsume mq message
	var delivery <-chan amqp.Delivery
	if delivery, err = c.channel.Consume(
		c.queue,       // queue
		c.consumerTag, // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	); err != nil {
		_ = c.channel.Close()
		_ = c.conn.Close()
		return err
	}

	go c.Handle(delivery)

	c.connNotify = c.conn.NotifyClose(make(chan *amqp.Error))
	c.channelNotify = c.channel.NotifyClose(make(chan *amqp.Error))

	return err
}

func (c *Consumer) ReConnect() {
	for {
		select {
		case err := <-c.connNotify:
			if err != nil {
				log.Println("rabbitmq consumer - connection NotifyClose: ", err)
			}
		case err := <-c.channelNotify:
			if err != nil {
				log.Println("rabbitmq consumer - channel NotifyClose: ", err)
			}
		case <-c.quit:
			return
		}

		// backstop
		if !c.conn.IsClosed() {
			// close message delivery
			if err := c.channel.Cancel(c.consumerTag, true); err != nil {
				log.Println("rabbitmq consumer - channel cancel failed: ", err)
			}

			if err := c.conn.Close(); err != nil {
				log.Println("rabbitmq consumer - channel cancel failed: ", err)
			}
		}

		// IMPORTANT: 必须清空 Notify，否则死连接不会释放
		for err := range c.channelNotify {
			println(err)
		}
		for err := range c.connNotify {
			println(err)
		}

	labelQuit:
		for {
			select {
			case <-c.quit:
				return
			default:
				if err := c.Run(); err != nil {
					log.Println("rabbitmq consumer - failCheck: ", err)

					// sleep 5s reconnect
					time.Sleep(time.Second * 5)
					continue
				}

				break labelQuit
			}
		}
	}
}

func (c *Consumer) Handle(delivery <-chan amqp.Delivery) {
	for d := range delivery {
		go func(delivery amqp.Delivery) {
			if err := c.handler(delivery.Body); err == nil {
				_ = delivery.Ack(false)
			} else {
				// 重新入队，否则未确认的消息会持续占用内存
				_ = delivery.Reject(true)
			}
		}(d)
	}
}
