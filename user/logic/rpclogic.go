package logic

import (
	"context"

	"integral-mall/user/model"
	"integral-mall/user/protos"
)

type (
	UserRpcServerLogic struct {
		userModel *model.UserModel
	}
)

func NewUserRpcServerLogic(userModel *model.UserModel) *UserRpcServerLogic {
	return &UserRpcServerLogic{userModel: userModel}
}

// FindUserByMobile(context.Context, *FindUserByMobileRequest) (*UserResponse, error)
//	FindId(context.Context, *FindIdRequest) (*UserResponse, error)
func (l *UserRpcServerLogic) FindUserByMobile(_ context.Context, r *protos.FindUserByMobileRequest) (*protos.UserResponse, error) {

	user, err := l.userModel.FindByMobile(r.Mobile)
	if err != nil {
		return nil, err
	}
	return &protos.UserResponse{
		Id: user.Id, Name: user.Name,
	}, nil
}
func (l *UserRpcServerLogic) FindId(_ context.Context, r *protos.FindIdRequest) (*protos.UserResponse, error) {
	user, err := l.userModel.FindById(r.Id)
	if err != nil {
		return nil, err
	}
	return &protos.UserResponse{
		Id: user.Id, Name: user.Name,
	}, nil
}
