package logic

import (
	"integral-mall/goods/model"
)

type (
	GoodsRpcServerLogic struct {
		goodsModel *model.GoodsModel
	}
)

func NewGoodsRpcServerLogic(goodsModel *model.GoodsModel) *GoodsRpcServerLogic {
	return &GoodsRpcServerLogic{goodsModel: goodsModel}
}

//
//func (l *UserRpcServerLogic) FindId(_ context.Context, r *protos.FindIdRequest) (*protos.UserResponse, error) {
//	user, err := l.userModel.FindById(r.Id)
//	if err != nil {
//		return nil, err
//	}
//	return &protos.UserResponse{
//		Id: user.Id, Name: user.Name,
//	}, nil
//}
