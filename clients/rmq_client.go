package clients

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	rabbitmqConnStr string
	channels        = make(map[string]*amqp.Channel)
)

type RmqClient interface {
	DeclareExAndQueue(exchange string, queue string)
	DeclareExAndQueueWithRKey(exchange string, queue string, key string)
	SendMsgToRMQ(exchange string, p any) (err error)
	ConsumeWithRetry(exchange string, queue string, processor MessageProcessor, secondSleep time.Duration)
	ConsumeWithKeyAndRetry(exchange string, queue string, routingKey string, processor MessageProcessor, secondSleep time.Duration)
}

type rmqClient struct {
	rmqConn *amqp.Connection
	user    string
	pass    string
	host    string
	appKey  string
	semafor sync.Mutex
}

func NewRmqClient(user string, pass string, host string, appKey string) *rmqClient {
	r := &rmqClient{
		user:    user,
		pass:    pass,
		host:    host,
		appKey:  appKey,
		semafor: sync.Mutex{},
	}
	r.init()
	return r
}

func (r *rmqClient) init() {
	rabbitmqConnStr = fmt.Sprintf("amqp://%s:%s@%s/", r.user, r.pass, r.host)
	r.rmqConn = r.connectToRabbitMQ()
}

func (r *rmqClient) connectToRabbitMQ() *amqp.Connection {
	conn, err := amqp.Dial(rabbitmqConnStr)
	if err != nil {
		logrus.Warn("Can't connect to RMQ", err)
	}
	return conn
}

func (r *rmqClient) DeclareExAndQueue(exchange string, queue string) {
	r.DeclareExAndQueueWithRKey(exchange, queue, "#")
}

func (r *rmqClient) DeclareExAndQueueWithRKey(exchange string, queue string, key string) {
	ch, err := r.getChannel()
	if err != nil {
		logrus.Warn("Can't get chanel from RMQ", err)
		return
	}
	err = r.DeclareTopicExchange(ch, exchange)
	if err != nil {
		logrus.Warn("Can't DeclareTopicExchange RMQ", err)
		return
	}
	q := r.DeclareQueue(queue, ch)
	err = ch.QueueBind(q.Name, key, exchange, false, nil)
	if err != nil {
		logrus.Warnf("Can't bind queue %s, exc: %s", q.Name, exchange)
		return
	}
	r.semafor.Lock()
	defer r.semafor.Unlock()
	channels[exchange] = ch
}

func (r *rmqClient) DeclareEx(exchange string) {
	ch, err := r.getChannel()
	if err != nil {
		logrus.Warn("Can't DeclareEx", err)
		return
	}
	err = r.DeclareTopicExchange(ch, exchange)
	if err != nil {
		logrus.Warn("Can't DeclareTopicExchange", err)
		return
	}
}

func (r *rmqClient) DeclareTopicExchange(ch *amqp.Channel, exchName string) error {
	err := ch.ExchangeDeclare(exchName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}
func (r *rmqClient) DeclareQueue(queue string, ch *amqp.Channel) (q amqp.Queue) {
	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logrus.Warn("Can't DeclareQueue", err)
		return
	}
	return q
}

func (r *rmqClient) SendMsgToRMQ(exchange string, p any) (err error) {
	return r.SendMsgToRMQWithKey(exchange, p, "")
}

func (r *rmqClient) SendMsgToRMQWithKey(exchange string, p any, key string) (err error) {
	contentType := "text/plain"
	return r.sendMsgToRMQWithKey(exchange, p, key, contentType)
}

func (r *rmqClient) sendMsgToRMQWithKey(exchange string, p any, key string, contentType string) (err error) {
	body, err := json.Marshal(p)
	if err != nil {
		logrus.Errorf("while marshal payload %s", err)
		return err
	}
	ch, err := r.getChannel()
	if err != nil {
		logrus.Errorf("can't get channel from connection %s ", err)
		return fmt.Errorf("can't get channel from connection %s ", err)
	}
	err = ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		},
	)
	defer ch.Close()

	if err != nil {
		logrus.Warn(err)
		return
	}
	return
}

func (r *rmqClient) SendJsonMsgToRMQWithKey(exchange string, p any, key string) (err error) {
	contentType := "application/json"
	return r.sendMsgToRMQWithKey(exchange, p, key, contentType)
}

func (r *rmqClient) getChannel() (*amqp.Channel, error) {
	if r.rmqConn.IsClosed() {
		r.Reconnect()
	}
	ch, err := r.rmqConn.Channel()
	if err != nil {
		logrus.Warn("Can't get chanel from RMQ", err)
		return nil, err
	}
	return ch, nil
}

func (r *rmqClient) GetMsgs(channelName string, qName string) <-chan amqp.Delivery {
	ch := channels[channelName]
	msgs, err := ch.Consume(
		qName, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		logrus.Warn(err)
		return nil
	}
	time.Sleep(5 * time.Second)
	return msgs
}

func (r *rmqClient) Reconnect() {
	r.rmqConn = r.connectToRabbitMQ()
}

type MessageProcessor func(msg amqp.Delivery) error

func (r *rmqClient) ConsumeWithKeyAndRetry(exchange string, queue string, routingKey string, processor MessageProcessor, secondSleep time.Duration) {
	for {
		r.DeclareExAndQueueWithRKey(exchange, queue, routingKey)
		msgs := r.GetMsgs(exchange, queue)
		if msgs == nil {
			logrus.Warn("Failed to get messages, retrying...")
			time.Sleep(30 * time.Second)
			continue
		}

		for msg := range msgs {
			err := processor(msg)
			if err != nil {
				logrus.Errorf("Error processing message: %v", err)
			}
			time.Sleep(secondSleep * time.Second)
		}

		logrus.Warn("Message channel closed, retrying...")
		time.Sleep(15 * time.Second)
	}
}

func (r *rmqClient) ConsumeWithRetry(exchange string, queue string, processor MessageProcessor, secondSleep time.Duration) {
	r.ConsumeWithKeyAndRetry(exchange, queue, "#", processor, secondSleep)
}
