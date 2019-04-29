package logic

import (
	"integral-mall/common/baseerror"
	"integral-mall/common/rpcxclient/integralrpcmodel"
	"integral-mall/goods/model"
)

type (
	GoodsLogic struct {
		goodsModel       *model.GoodsModel
		integralRpcModel *integralrpcmodel.IntegralRpcModel
	}
	GoodSearchRequest struct {
		Name string `json:"name"`
		Page int    `json:"page"`
	}

	GoodsOrderRequest struct {
		Id  int64 `json:"id"  binding:"required"`
		Num int64 `json:"num" binding:"required"`
	}
	GoodsOrderResponse struct {
	}

	GoodSearchResponse struct {
		Total     int64        `json:"total"`
		GoodsList []*GoodsView `json:"goodsList"`
	}

	GoodsView struct {
		Name  string `json:"name"`
		Id    int64  `json:"id"`
		Image string `json:"image"`
		Intro string `json:"intro"`
		Price int    `json:"price"`
		Store int64  `json:"store"`
	}
)

var (
	ErrStoreOver    = baseerror.NewBaseError("商品库存没了")
	ErrIntegralOver = baseerror.NewBaseError("您的积分不能买这个!!!!")
)

func NewGoodsLogic(goodsModel *model.GoodsModel,
	integralRpcModel *integralrpcmodel.IntegralRpcModel,
) *GoodsLogic {
	return &GoodsLogic{goodsModel: goodsModel, integralRpcModel: integralRpcModel}
}

func (l *GoodsLogic) GoodSearch(r *GoodSearchRequest) (*GoodSearchResponse, error) {
	goodsList, count, err := l.goodsModel.PageList(r.Name, r.Page)
	if err != nil {
		return nil, err
	}
	response := &GoodSearchResponse{Total: count, GoodsList: []*GoodsView(nil)}
	for _, g := range goodsList {
		response.GoodsList = append(response.GoodsList, &GoodsView{
			Id: g.Id, Name: g.GoodName, Image: g.Image, Store: g.Store, Intro: g.Intro,
		})
	}
	return response, nil
}

func (l *GoodsLogic) GoodsOrder(r *GoodsOrderRequest, userId int) (*GoodsOrderResponse, error) {
	goods, err := l.goodsModel.FindById(r.Id)
	if err != nil {
		return nil, err
	}
	if goods.Store <= 0 {
		return nil, ErrStoreOver
	}
	integral, err := l.integralRpcModel.FindOneByUserId(userId)
	if err != nil {
		return nil, err
	}
	if integral.Integral < goods.Price {
		return nil, ErrIntegralOver
	}
	if err := l.goodsModel.TransactionChangeStore(r.Id, r.Num, userId, func(userId int) error {
		//integral
		if err := l.integralRpcModel.ConsumerIntegral(userId, goods.Price); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return new(GoodsOrderResponse), nil
}
