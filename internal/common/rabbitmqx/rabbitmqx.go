package rabbitmqx

import (
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Service struct {
	sync.RWMutex
	AmqpUrl       string //amqp地址
	ConnectionNum int    //连接数
	ChannelNum    int    //每个连接的channel数量

	connections  map[int]*connection
	channels     map[int]*channel
	idelChannels []int
	busyChannels map[int]int
}

type channel struct {
	channelId   int
	ch          *amqp.Channel
	notifyClose chan *amqp.Error
}

type connection struct {
	connectId   int
	conn        *amqp.Connection
	notifyClose chan *amqp.Error
}

var amqpServer *Service

func GetMq() *Service {
	return amqpServer
}

func Init(_host, _port, _user, _password string, _connectionNum, _channelNum int) error {
	amqpUrl := "amqp://" + _user + ":" + _password + "@" + _host + ":" + _port + "/"
	amqpServer = &Service{
		AmqpUrl:       amqpUrl,
		ConnectionNum: _connectionNum,
		ChannelNum:    _channelNum,
	}

	rand.Seed(time.Now().UnixNano())
	if amqpServer.ConnectionNum == 0 {
		amqpServer.ConnectionNum = 10
	}
	if amqpServer.ChannelNum == 0 {
		amqpServer.ChannelNum = 10
	}
	amqpServer.busyChannels = make(map[int]int)
	if err := amqpServer.connectPool(); err != nil {
		return err
	}
	if err := amqpServer.channelPool(); err != nil {
		return err
	}
	return nil
}

func (S *Service) connectPool() error {
	S.connections = make(map[int]*connection)
	for i := 0; i < S.ConnectionNum; i++ {
		connection, err := S.connect(i)
		if err != nil {
			return err
		}
		S.connections[i] = connection
	}
	return nil
}

func (S *Service) channelPool() error {
	S.channels = make(map[int]*channel)
	for i, connection := range S.connections {
		for j := 0; j < S.ChannelNum; j++ {
			channelId := i*S.ChannelNum + j
			channel, err := S.createChannel(channelId, connection)
			if err != nil {
				return err
			}
			S.channels[channelId] = channel
			S.idelChannels = append(S.idelChannels, channelId)
		}
	}
	return nil
}

func (S *Service) connect(_connectId int) (*connection, error) {
	var notifyClose = make(chan *amqp.Error)

	connection := &connection{
		notifyClose: notifyClose,
		connectId:   _connectId,
	}
	conn, err := amqp.Dial(S.AmqpUrl)
	if err != nil {
		return connection, err
	}
	connection.conn = conn
	conn.NotifyClose(connection.notifyClose)
	go func() {
		select {
		case <-connection.notifyClose:
			log.Println("close connectId:", connection.connectId)
		}
	}()
	return connection, err
}

func (S *Service) createChannel(_channelId int, _connection *connection) (*channel, error) {
	var notifyClose = make(chan *amqp.Error)

	cha := &channel{
		notifyClose: notifyClose,
		channelId:   _channelId,
	}
	ch, err := _connection.conn.Channel()
	if err != nil {
		return cha, err
	}
	cha.ch = ch
	ch.NotifyClose(cha.notifyClose)
	go func() {
		select {
		case <-cha.notifyClose:
			log.Println("mq close channelId:", cha.channelId)
		}
	}()
	return cha, nil
}

func (S *Service) getChannel() (*amqp.Channel, int) {
	S.Lock()
	defer S.Unlock()
	idelLength := len(S.idelChannels)
	if idelLength > 0 {
		index := rand.Intn(idelLength)
		channelId := S.idelChannels[index]
		S.idelChannels = append(S.idelChannels[:index], S.idelChannels[index+1:]...)
		S.busyChannels[channelId] = channelId

		ch := S.channels[channelId].ch
		return ch, channelId
	} else {
		return nil, -1
	}
}

func (S *Service) lockWriteConnect(_connectId int, _newConn *connection) {
	S.Lock()
	defer S.Unlock()
	S.connections[_connectId] = _newConn
}

func (S *Service) lockWriteChannel(_channelId int, _cha *channel) {
	S.Lock()
	defer S.Unlock()
	S.channels[_channelId] = _cha
}

func (S *Service) backChannelId(_channelId int) {
	S.Lock()
	defer S.Unlock()
	S.idelChannels = append(S.idelChannels, _channelId)
	delete(S.busyChannels, _channelId)
	return
}

func (S *Service) publish(ch *amqp.Channel, exchangeName string, routeKey string, _data []byte) (err error) {
	err = ch.Publish(
		exchangeName, // exchange
		routeKey,     // routing key
		false,        // 设置为true时，至少将该消息route到一个队列中，否则返还给生产者；false时，上述情形直接丢弃！
		false,        // 设置为true时，queue上有消费者，马上投递，没有消费者，返还生产者不进队列！
		amqp.Publishing{
			ContentType: "application/json",
			Body:        _data,
		})
	return
}

func (S *Service) reconnect(_channelId int) error {
	S.Lock()
	defer S.Unlock()
	connectId := int(_channelId / S.ChannelNum)
	connection := S.connections[connectId]
	if connection.conn.IsClosed() {
		connection, err := S.connect(connectId)
		log.Printf("reconnect connectId:%v, err:%v", connectId, err)
		if err != nil {
			return err
		}
		S.connections[connectId] = connection

		for j := 0; j < S.ChannelNum; j++ {
			channelId := connectId*S.ChannelNum + j
			err2 := S.channels[channelId].ch.Close()
			channel, err := S.createChannel(connectId, connection)
			log.Printf("reconnect connectId:%v_channelId:%v, err:%v, err2:%v", connectId, channelId, err, err2)
			if err != nil {
				return err
			}
			S.channels[channelId] = channel
		}
		return nil
	}

	err2 := S.channels[_channelId].ch.Close()
	channel, err := S.createChannel(connectId, connection)
	log.Printf("reconnect _channelId:%v, err:%v, err2:%v", _channelId, err, err2)
	if err != nil {
		return err
	}
	S.channels[_channelId] = channel
	return nil
}

func (S *Service) PutIntoQueue(_exchangeName string, _routeKey string, _data []byte) (puberr error) {
	ch, channelId := S.getChannel()
	defer func() {
		if err := recover(); err != nil {
			puberrMsg, _ := err.(string)
			puberr = errors.Errorf(puberrMsg)
		}
	}()
	defer func() {
		if channelId >= 0 {
			if puberr != nil && strings.Index(puberr.Error(), "channel/connection is not open") >= 0 {
				if err := S.reconnect(channelId); err == nil {
					S.publish(S.channels[channelId].ch, _exchangeName, _routeKey, _data)
				}
			}
			S.backChannelId(channelId)
		}
	}()

	if ch == nil {
		connectId := rand.Intn(S.ConnectionNum)
		cha, err := S.createChannel(connectId, S.connections[connectId])
		if err != nil {
			return err
		}
		defer cha.ch.Close()
		ch = cha.ch
	}

	puberr = S.publish(ch, _exchangeName, _routeKey, _data)
	return
}

func (S *Service) Consume(_queueName string, _autoAck bool, _exclusive bool, _noLocal bool, _noWait bool, _handler func([]byte) error) {
	for {
		ch, channelId := S.getChannel()
		if ch == nil {
			continue
		}

		delivery, err := ch.Consume(_queueName, "", _autoAck, _exclusive, _noLocal, _noWait, nil)
		if err == nil {
			for d := range delivery {
				err := _handler(d.Body)
				if !_autoAck {
					if err == nil {
						_ = d.Ack(true)
					} else {
						// 重新入队，否则未确认的消息会持续占用内存
						_ = d.Reject(true)
					}
				}
			}
		}
		S.reconnect(channelId)
		S.backChannelId(channelId)
		time.Sleep(time.Second)
	}
}

func (S *Service) ExchangeDeclare(_exchangeName string, _exchangeKind string, _durable, _autoDelete, _internal, _noWait bool) (puberr error) {
	ch, channelId := S.getChannel()
	defer func() {
		if channelId >= 0 {
			S.backChannelId(channelId)
		}
	}()
	if ch == nil {
		connectId := rand.Intn(S.ConnectionNum)
		cha, err := S.createChannel(connectId, S.connections[connectId])
		if err != nil {
			return err
		}
		defer cha.ch.Close()
		ch = cha.ch
	}

	puberr = ch.ExchangeDeclare(_exchangeName, _exchangeKind, _durable, _autoDelete, _internal, _noWait, nil)
	return
}

func (S *Service) QueueDeclare(_queueName string, _durable, _autoDelete, _exclusive, _noWait bool) (puberr error) {
	ch, channelId := S.getChannel()
	defer func() {
		if channelId >= 0 {
			S.backChannelId(channelId)
		}
	}()
	if ch == nil {
		connectId := rand.Intn(S.ConnectionNum)
		cha, err := S.createChannel(connectId, S.connections[connectId])
		if err != nil {
			return err
		}
		defer cha.ch.Close()
		ch = cha.ch
	}

	_, puberr = ch.QueueDeclare(_queueName, _durable, _autoDelete, _exclusive, _noWait, nil)
	return
}

func (S *Service) QueueBind(_exchangeName string, _queueName string, _bindKey string, _noWait bool) (puberr error) {
	ch, channelId := S.getChannel()
	defer func() {
		if channelId >= 0 {
			S.backChannelId(channelId)
		}
	}()
	if ch == nil {
		connectId := rand.Intn(S.ConnectionNum)
		cha, err := S.createChannel(connectId, S.connections[connectId])
		if err != nil {
			return err
		}
		defer cha.ch.Close()
		ch = cha.ch
	}

	puberr = ch.QueueBind(_queueName, _bindKey, _exchangeName, _noWait, nil)
	return
}

func (S *Service) GetQueueMsgCount(_queueName string) (int, error) {
	count := 0
	ch, channelId := S.getChannel()
	defer func() {
		if channelId >= 0 {
			S.backChannelId(channelId)
		}
	}()
	if ch == nil {
		connectId := rand.Intn(S.ConnectionNum)
		cha, err := S.createChannel(connectId, S.connections[connectId])
		if err != nil {
			return count, err
		}
		defer cha.ch.Close()
		ch = cha.ch
	}
	queue, err := ch.QueueInspect(_queueName)
	if err != nil {
		return count, err
	}
	return queue.Messages, nil
}

func (S *Service) QueuePurage(_queueName string) error {
	ch, channelId := S.getChannel()
	defer func() {
		if channelId >= 0 {
			S.backChannelId(channelId)
		}
	}()
	if ch == nil {
		connectId := rand.Intn(S.ConnectionNum)
		cha, err := S.createChannel(connectId, S.connections[connectId])
		if err != nil {
			return err
		}
		defer cha.ch.Close()
		ch = cha.ch
	}
	_, err := ch.QueuePurge(_queueName, false)
	return err
}
