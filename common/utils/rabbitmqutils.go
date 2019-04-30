package utils

import (
	"github.com/streadway/amqp"
	"github.com/yakaa/log4g"
)

type (
	RabbitMqServer struct {
		dialHost     string
		queueName    string
		rabbitMqConn *amqp.Connection
		channel      *amqp.Channel
	}
)

//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
//failOnError(err, "Failed to connect to RabbitMQ")
func NewRabbitMqServer(dialHost, queueName string) (*RabbitMqServer, error) {
	rabbitMqServer := &RabbitMqServer{dialHost: dialHost, queueName: queueName}
	if err := rabbitMqServer.createDial(); err != nil {
		return nil, err
	}
	return rabbitMqServer, nil
}

func (l *RabbitMqServer) createDial() error {
	conn, err := amqp.Dial(l.dialHost)
	if err != nil {
		return err
	}
	l.rabbitMqConn = conn
	l.channel, err = l.rabbitMqConn.Channel()
	if err != nil {
		return nil
	}
	return nil
}

func (l *RabbitMqServer) CloseRabbitMqConn() {
	if err := l.rabbitMqConn.Close(); err != nil {
		log4g.ErrorFormat("CloseRabbitMqConn err %+v", err)
	}
	if l.channel != nil {
		if err := l.channel.Close(); err != nil {
			log4g.ErrorFormat("ConsumeChannel err %+v", err)
		}
	}
}

func (l *RabbitMqServer) PushMessage(message string) {
	q, err := l.QueueDeclare(l.channel)
	if err != nil {
		log4g.ErrorFormat("PushMessage err %+v", err)
		return
	}
	err = l.channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{Body: []byte(message)})
	if err != nil {
		log4g.ErrorFormat("ch.Publish err %+v", err)
		return
	}
}

func (l *RabbitMqServer) ConsumeMessage(consumeMessageFunc func(message string) error) {
	q, err := l.QueueDeclare(l.channel)
	if err != nil {
		log4g.ErrorFormat("PushMessage err %+v", err)
		return
	}
	messageList, err := l.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log4g.ErrorFormat("ch.Consume err %+v", err)
		return
	}
	go func() {
		for d := range messageList {
			msg := string(d.Body)
			log4g.InfoFormat("get message %s", msg)
			if err := consumeMessageFunc(msg); err != nil {
				l.PushMessage(msg)
			} else {
				log4g.InfoFormat("Consume message %s [SUCESSS]", msg)
			}
		}
	}()
}

func (l *RabbitMqServer) QueueDeclare(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		l.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
}
