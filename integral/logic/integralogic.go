package logic

import (
	"context"

	"github.com/streadway/amqp"
	"github.com/yakaa/log4g"

	"integral-mall/integral/model"
	"integral-mall/integral/protos"
)

type (
	IntegralLogic struct {
		dialHost      string
		queueName     string
		integralModel *model.IntegralModel
		rabbitMqConn  *amqp.Connection
		channel       *amqp.Channel
	}
)

//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
//failOnError(err, "Failed to connect to RabbitMQ")
func NewIntegralLogic(dialHost, queueName string, integralModel *model.IntegralModel) (*IntegralLogic, error) {
	integralLogic := &IntegralLogic{dialHost: dialHost, queueName: queueName, integralModel: integralModel}
	if err := integralLogic.createDial(); err != nil {
		return nil, err
	}
	return integralLogic, nil
}

func (l *IntegralLogic) createDial() error {
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

func (l *IntegralLogic) CloseRabbitMqConn() {
	if err := l.rabbitMqConn.Close(); err != nil {
		log4g.ErrorFormat("CloseRabbitMqConn err %+v", err)
	}
	if l.channel != nil {
		if err := l.channel.Close(); err != nil {
			log4g.ErrorFormat("ConsumeChannel err %+v", err)
		}
	}
}

func (l *IntegralLogic) PushMessage(message string) {
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

func (l *IntegralLogic) ConsumeMessage() {
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
			if err := l.integralModel.ExecSql(msg); err != nil {
				l.PushMessage(msg)
			} else {
				log4g.InfoFormat("Consume message %s [SUCESSS]", msg)
			}
		}
	}()
}

func (l *IntegralLogic) QueueDeclare(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		l.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
}

//AddIntegral(context.Context, *AddIntegralRequest) (*IntegralResponse, error)
//	ConsumerIntegral(context.Context, *ConsumerIntegralRequest) (*IntegralResponse, error)
func (l *IntegralLogic) AddIntegral(_ context.Context, r *protos.AddIntegralRequest) (*protos.IntegralResponse, error) {
	l.PushMessage(l.integralModel.InsertIntegralSql(int(r.UserId), int(r.Integral)))
	return &protos.IntegralResponse{
		UserId: r.UserId, Integral: r.Integral,
	}, nil
}

func (l *IntegralLogic) ConsumerIntegral(_ context.Context, r *protos.ConsumerIntegralRequest) (*protos.IntegralResponse, error) {
	l.PushMessage(l.integralModel.UpdateIntegralByUserIdSql(int(r.UserId), int(r.ConsumerIntegral)))
	return new(protos.IntegralResponse), nil
}

func (l *IntegralLogic) FindOneByUserId(_ context.Context, r *protos.FindOneByUserIdRequest) (*protos.IntegralResponse, error) {
	//l.PushMessage(l.integralModel.UpdateIntegralByUserIdSql(int(r.UserId), int(r.ConsumerIntegral)))
	one, err := l.integralModel.FindByUserId(int(r.UserId))
	if err != nil {
		return nil, err
	}
	return &protos.IntegralResponse{
		UserId: r.UserId, Integral: int64(one.Integral),
	}, nil
}
