package orderrpcmodel

import (
	"context"

	"github.com/yakaa/grpcx"
	"github.com/yakaa/grpcx/config"

	"integral-mall/order/protos"
)

//r, err := grpcx.MustNewGrpcxClient(conf)
//	if err != nil {
//		panic(err)
//	}

type (
	OrderModel struct {
		cli *grpcx.GrpcxClient
	}
	OrderBookingResponse struct {
		Success bool
	}
	OrderOneResponse struct {
		Id         string `json:"id,omitempty"`
		GoodName   string `json:"goodName,omitempty"`
		Mobile     string `json:"mobile,omitempty"`
		Num        int64  `json:"num,omitempty"`
		UserId     int64  `json:"userId,omitempty"`
		CreateTime string `json:"createTime,omitempty"`
	}
)

func NewOrderModel(cli *grpcx.GrpcxClient) *OrderModel {

	return &OrderModel{cli: cli}
}

func (m *OrderModel) BookingGoods(orderId, goodsName, mobile string, goodsId, userId, num int64) (*OrderBookingResponse, error) {
	conn, err := m.cli.GetConnection()
	if err != nil {
		return nil, err
	}
	client := protos.NewOrderRpcClient(conn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), config.GrpcxDialTimeout)
	defer cancelFunc()
	_, err = client.BookingGoods(
		ctx,
		&protos.BookingGoodsRequest{
			OrderId:   orderId,
			GoodsName: goodsName,
			GoodsId:   goodsId,
			UserId:    userId,
			Num:       num,
			Mobile:    mobile,
		})
	if err != nil {
		return nil, err
	}
	return &OrderBookingResponse{Success: true}, nil
}

func (m *OrderModel) FindId(orderId string) (*OrderOneResponse, error) {
	conn, err := m.cli.GetConnection()
	if err != nil {
		return nil, err
	}
	client := protos.NewOrderRpcClient(conn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), config.GrpcxDialTimeout)
	defer cancelFunc()
	order, err := client.FindOrderId(
		ctx,
		&protos.FindOrderIdRequest{OrderId: orderId})
	if err != nil {
		return nil, err
	}
	return &OrderOneResponse{
		Id:         order.Id,
		GoodName:   order.GoodName,
		Mobile:     order.Mobile,
		Num:        int64(order.Num),
		UserId:     order.UserId,
		CreateTime: order.CreateTime,
	}, nil
}
