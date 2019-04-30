package controller

import (
	"github.com/gin-gonic/gin"

	"integral-mall/common/baseresponse"
	"integral-mall/goods/logic"
)

type (
	GoodsController struct {
		goodsLogic *logic.GoodsLogic
	}
)

func NewGoodsController(goodsLogic *logic.GoodsLogic) *GoodsController {
	return &GoodsController{goodsLogic: goodsLogic}
}

func (c *GoodsController) GoodSearch(ctx *gin.Context) {
	r := new(logic.GoodSearchRequest)
	if err := ctx.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(ctx, err)
		return
	}
	res, err := c.goodsLogic.GoodSearch(r)
	baseresponse.HttpResponse(ctx, res, err)
	return
}

func (c *GoodsController) GoodsOrder(ctx *gin.Context) {
	r := new(logic.GoodsOrderRequest)
	if err := ctx.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(ctx, err)
		return
	}
	r.UserId = ctx.GetInt("userId")
	res, err := c.goodsLogic.GoodsOrder(r)
	baseresponse.HttpResponse(ctx, res, err)
	return
}
