package logic

import (
	"integral-mall/common/baseerror"
	"integral-mall/common/rpcxclient/integralrpcmodel"
	"integral-mall/common/rpcxclient/orderrpcmodel"
	"integral-mall/goods/model"

	"github.com/satori/go.uuid"
)

type (
	GoodsLogic struct {
		goodsModel       *model.GoodsModel
		integralRpcModel *integralrpcmodel.IntegralRpcModel
		orderRpcModel    *orderrpcmodel.OrderModel
	}
	GoodSearchRequest struct {
		Name string `json:"name"`
		Page int    `json:"page"`
	}

	GoodsOrderRequest struct {
		Id     int64  `json:"id"  binding:"required"`
		UserId int    `json:"-"`
		Num    int64  `json:"num" binding:"required"`
		Mobile string `json:"mobile" binding:"required"`
	}
	GoodsOrderResponse struct {
		OrderId string `json:"orderId"`
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
	orderRpcModel *orderrpcmodel.OrderModel,
) *GoodsLogic {
	return &GoodsLogic{goodsModel: goodsModel, integralRpcModel: integralRpcModel, orderRpcModel: orderRpcModel}
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

func (l *GoodsLogic) GoodsOrder(r *GoodsOrderRequest) (*GoodsOrderResponse, error) {
	goods, err := l.goodsModel.FindById(r.Id)
	if err != nil {
		return nil, err
	}
	if goods.Store <= 0 {
		return nil, ErrStoreOver
	}
	integral, err := l.integralRpcModel.FindOneByUserId(r.UserId)
	if err != nil {
		return nil, err
	}
	if integral.Integral < goods.Price {
		return nil, ErrIntegralOver
	}
	orderId := uuid.NewV4().String()
	if err := l.goodsModel.TransactionChangeStore(r.Id, r.Num, r.UserId, func(userId int) error {
		//integral
		if err := l.integralRpcModel.ConsumerIntegral(userId, goods.Price); err != nil {
			return err
		}
		// order
		if _, err := l.orderRpcModel.BookingGoods(
			orderId, goods.GoodName, r.Mobile, r.Id, int64(r.UserId), r.Num,
		); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return new(GoodsOrderResponse), nil
}
