package integralrpcmodel

import (
	"context"

	"github.com/yakaa/grpcx"
	"github.com/yakaa/grpcx/config"

	"integral-mall/integral/protos"
)

//r, err := grpcx.MustNewGrpcxClient(conf)
//	if err != nil {
//		panic(err)
//	}

type (
	IntegralRpcModel struct {
		cli *grpcx.GrpcxClient
	}
)

func NewIntegralRpcModel(cli *grpcx.GrpcxClient) *IntegralRpcModel {

	return &IntegralRpcModel{cli: cli}
}

func (m *IntegralRpcModel) AddIntegral(userId, integral int) error {
	conn, err := m.cli.GetConnection()
	if err != nil {
		return err
	}
	clientIntegral := protos.NewIntegralRpcClient(conn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), config.GrpcxDialTimeout)
	defer cancelFunc()
	if _, err := clientIntegral.AddIntegral(
		ctx,
		&protos.AddIntegralRequest{UserId: int64(userId),
			Integral: int64(integral)}); err != nil {
		return err
	}
	return nil
}

func (m *IntegralRpcModel) ConsumerIntegral(userId, integral int) error {
	conn, err := m.cli.GetConnection()
	if err != nil {
		return err
	}
	clientIntegral := protos.NewIntegralRpcClient(conn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), config.GrpcxDialTimeout)
	defer cancelFunc()
	if _, err := clientIntegral.ConsumerIntegral(
		ctx,
		&protos.ConsumerIntegralRequest{UserId: int64(userId),
			ConsumerIntegral: int64(integral)}); err != nil {
		return err
	}
	return nil
}
