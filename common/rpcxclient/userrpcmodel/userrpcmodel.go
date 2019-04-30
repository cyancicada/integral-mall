package userrpcmodel

import (
	"context"

	"github.com/yakaa/grpcx"
	"github.com/yakaa/grpcx/config"

	"integral-mall/user/protos"
)

//r, err := grpcx.MustNewGrpcxClient(conf)
//	if err != nil {
//		panic(err)
//	}

type (
	UserRpcModel struct {
		cli *grpcx.GrpcxClient
	}
	UserRpcClientModel struct {
		UserId   int64
		UserName string
	}
)

func NewUserRpcModel(cli *grpcx.GrpcxClient) *UserRpcModel {

	return &UserRpcModel{cli: cli}
}

func (m *UserRpcModel) FindUserByMobile(mobile string) (*UserRpcClientModel, error) {
	conn, err := m.cli.GetConnection()
	if err != nil {
		return nil, err
	}
	client := protos.NewUserRpcClient(conn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), config.GrpcxDialTimeout)
	defer cancelFunc()
	res, err := client.FindUserByMobile(
		ctx,
		&protos.FindUserByMobileRequest{Mobile: mobile})
	if err != nil {
		return nil, err
	}
	return &UserRpcClientModel{UserName: res.Name, UserId: res.Id}, nil
}

func (m *UserRpcModel) FindId(userId int) (*UserRpcClientModel, error) {
	conn, err := m.cli.GetConnection()
	if err != nil {
		return nil, err
	}
	client := protos.NewUserRpcClient(conn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), config.GrpcxDialTimeout)
	defer cancelFunc()
	res, err := client.FindId(
		ctx,
		&protos.FindIdRequest{Id: int64(userId)})
	if err != nil {
		return nil, err
	}
	return &UserRpcClientModel{UserName: res.Name, UserId: res.Id}, nil
}
