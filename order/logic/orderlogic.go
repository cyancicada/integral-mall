package logic

import (
	"github.com/go-redis/redis"

	"integral-mall/common/baseerror"
	"integral-mall/common/i18n"
	"integral-mall/common/rpcxclient/userrpcmodel"
	"integral-mall/order/model"
)

type (
	OrderLogic struct {
		orderModel   *model.OrderModel
		redisCache   *redis.Client
		userRpcModel *userrpcmodel.UserRpcModel
	}
	OrderListRequest struct {
		UserId int64 `json:"-"`
		Page   int   `json:"page"`
	}
	OrderListResponse struct {
		Total    int64            `json:"total"`
		DataList []*OrderListView `json:"dataList"`
	}
	OrderListView struct {
		Id         string `json:"id"`
		GoodsName  string `json:"goodsName"`
		Mobile     string `json:"mobile"`
		Num        int64  `json:"num"`
		UserId     int64  `json:"userId"`
		UserName   string `json:"userName"`
		CreateTime string `json:"createTime"`
	}
)

var (
	ErrUserNotFound = baseerror.NewBaseError("用户不存在")
)

func NewOrderLogic(orderModel *model.OrderModel,
	redisCache *redis.Client,
	userRpcModel *userrpcmodel.UserRpcModel,
) *OrderLogic {

	return &OrderLogic{orderModel: orderModel, redisCache: redisCache, userRpcModel: userRpcModel}
}

func (l *OrderLogic) OrderList(r *OrderListRequest) (*OrderListResponse, error) {
	orders, count, err := l.orderModel.PageFindByUserId(r.UserId, r.Page)
	if err != nil {
		return nil, err
	}
	response := &OrderListResponse{Total: count, DataList: []*OrderListView(nil)}
	userInfo, err := l.userRpcModel.FindId(int(r.UserId))
	if err != nil {
		return nil, ErrUserNotFound
	}
	for _, order := range orders {
		response.DataList = append(response.DataList, &OrderListView{
			Id:         order.Id,
			GoodsName:  order.GoodName,
			Mobile:     order.Mobile,
			Num:        order.Num,
			UserId:     r.UserId,
			UserName:   userInfo.UserName,
			CreateTime: order.CreateTime.Format(i18n.TimLayOut),
		})
	}
	return response, nil
}
