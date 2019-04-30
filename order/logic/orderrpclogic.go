package logic

import (
	"context"

	"github.com/yakaa/log4g"

	"integral-mall/common/i18n"
	"integral-mall/common/utils"
	"integral-mall/order/model"
	"integral-mall/order/protos"
)

type (
	OrderRpcServerLogic struct {
		orderModel     *model.OrderModel
		rabbitMqServer *utils.RabbitMqServer
	}
)

func NewOrderRpcServerLogic(orderModel *model.OrderModel,

	rabbitMqServer *utils.RabbitMqServer,
) *OrderRpcServerLogic {
	return &OrderRpcServerLogic{orderModel: orderModel, rabbitMqServer: rabbitMqServer}
}

func (l *OrderRpcServerLogic) ConsumeMessageStart() {
	l.rabbitMqServer.ConsumeMessage(func(message string) error {
		log4g.InfoFormat("ConsumeMessageStart message %s", message)
		return l.orderModel.ExecSql(message)
	})
}

func (l *OrderRpcServerLogic) CloseRabbitMqServer() {
	l.rabbitMqServer.CloseRabbitMqConn()
}

//BookingGoods(ctx context.Context, in *BookingGoodsRequest, opts ...grpc.CallOption) (*BookingGoodsRequestResponse, error)
//	FindId(ctx context.Context, in *FindIdRequest, opts ...grpc.CallOption) (*OrderOneResponse, error)
func (l *OrderRpcServerLogic) BookingGoods(_ context.Context, r *protos.BookingGoodsRequest) (*protos.BookingGoodsResponse, error) {

	l.rabbitMqServer.PushMessage(
		l.orderModel.BookingGoodsSql(r.OrderId, r.GoodsId, r.GoodsName, r.Mobile, r.UserId, r.Num),
	)
	return &protos.BookingGoodsResponse{
		OrderId: r.OrderId,
	}, nil
}

func (l *OrderRpcServerLogic) FindOrderId(_ context.Context, r *protos.FindOrderIdRequest) (*protos.OrderOneResponse, error) {
	order, err := l.orderModel.FindById(r.OrderId)
	if err != nil {
		return nil, err
	}
	return &protos.OrderOneResponse{
		Id:         order.Id,
		GoodName:   order.GoodName,
		Mobile:     order.Mobile,
		Num:        int64(order.Num),
		UserId:     order.UserId,
		CreateTime: order.CreateTime.Format(i18n.TimLayOut),
	}, nil
}
