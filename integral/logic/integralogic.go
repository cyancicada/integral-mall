package logic

import (
	"context"

	"integral-mall/common/utils"
	"integral-mall/integral/model"
	"integral-mall/integral/protos"
)

type (
	IntegralLogic struct {
		rabbitMqServer *utils.RabbitMqServer
		integralModel  *model.IntegralModel
	}
)

//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
//failOnError(err, "Failed to connect to RabbitMQ")
func NewIntegralLogic(rabbitMqServer *utils.RabbitMqServer, integralModel *model.IntegralModel) *IntegralLogic {

	return &IntegralLogic{rabbitMqServer: rabbitMqServer, integralModel: integralModel}
}

func (l *IntegralLogic) ConsumeMessage() {
	l.rabbitMqServer.ConsumeMessage(func(message string) error {
		return l.integralModel.ExecSql(message)
	})
}

func (l *IntegralLogic) Close() {
	l.rabbitMqServer.CloseRabbitMqConn()
}

//AddIntegral(context.Context, *AddIntegralRequest) (*IntegralResponse, error)
//	ConsumerIntegral(context.Context, *ConsumerIntegralRequest) (*IntegralResponse, error)
func (l *IntegralLogic) AddIntegral(_ context.Context, r *protos.AddIntegralRequest) (*protos.IntegralResponse, error) {
	err := l.rabbitMqServer.PushMessage(l.integralModel.InsertIntegralSql(int(r.UserId), int(r.Integral)))
	return &protos.IntegralResponse{
		UserId: r.UserId, Integral: r.Integral,
	}, err
}

func (l *IntegralLogic) ConsumerIntegral(_ context.Context, r *protos.ConsumerIntegralRequest) (*protos.IntegralResponse, error) {
	err := l.rabbitMqServer.PushMessage(l.integralModel.UpdateIntegralByUserIdSql(int(r.UserId), int(r.ConsumerIntegral)))
	return new(protos.IntegralResponse), err
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
