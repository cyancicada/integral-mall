package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"integral-mall/common/baseerror"
	"integral-mall/common/baseresponse"
)

type (
	Authorization struct {
		redisCache *redis.Client
	}
)

var (
	ErrAuthorization = baseerror.NewBaseError("请先登录")
)

func NewAuthorization(redisCache *redis.Client) *Authorization {
	return &Authorization{redisCache: redisCache}
}
func (a *Authorization) Auth(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	if strings.TrimSpace(header) == "" {
		baseresponse.HttpResponse(ctx, nil, ErrAuthorization)
		ctx.Abort()
	}
	sc := a.redisCache.Get(header)
	userId, _ := strconv.Atoi(sc.Val())
	ctx.Set("userId", userId)
	ctx.Next()
}
