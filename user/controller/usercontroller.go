package controller

import (
	"github.com/gin-gonic/gin"

	"integral-mall/common/baseresponse"
	"integral-mall/user/logic"
)

type (
	UserController struct {
		userLogic *logic.UserLogic
	}
)

func NewUserController(userLogic *logic.UserLogic) *UserController {

	return &UserController{userLogic: userLogic}
}

//register

func (c *UserController) Register(ctx *gin.Context) {
	r := new(logic.RegisterRequest)
	if err := ctx.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(ctx, err)
		return
	}
	res, err := c.userLogic.Register(r)
	baseresponse.HttpResponse(ctx, res, err)
	return
}

//Login

func (c *UserController) Login(ctx *gin.Context) {
	r := new(logic.LoginRequest)
	if err := ctx.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(ctx, err)
		return
	}
	res, err := c.userLogic.Login(r)
	baseresponse.HttpResponse(ctx, res, err)
	return
}
