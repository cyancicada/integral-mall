package controller

import (
	"github.com/gin-gonic/gin"

	"integral-mall/common/baseresponse"
	"integral-mall/order/logic"
)

type (
	OrderController struct {
		orderLogic *logic.OrderLogic
	}
)

func NewOrderController(orderLogic *logic.OrderLogic) *OrderController {

	return &OrderController{orderLogic: orderLogic}
}

//register

func (c *OrderController) OrderList(ctx *gin.Context) {
	r := new(logic.OrderListRequest)
	if err := ctx.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(ctx, err)
		return
	}
	r.UserId = int64(ctx.GetInt("userId"))
	res, err := c.orderLogic.OrderList(r)
	baseresponse.HttpResponse(ctx, res, err)
	return
}
