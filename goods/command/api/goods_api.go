package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"github.com/yakaa/grpcx"
	"github.com/yakaa/log4g"

	_ "github.com/go-sql-driver/mysql"

	"integral-mall/common/middleware"
	"integral-mall/common/rpcxclient/integralrpcmodel"
	"integral-mall/goods/command/api/config"
	"integral-mall/goods/controller"
	"integral-mall/goods/logic"
	"integral-mall/goods/model"
)

var configFile = flag.String("f", "config/config.json", "use config")

func main() {
	flag.Parse()
	conf := new(config.Config)
	bs, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bs, conf); err != nil {
		log.Fatal(err)
	}
	log4g.Init(log4g.Config{Path: "logs"})
	gin.DefaultWriter = log4g.InfoLog
	gin.DefaultErrorWriter = log4g.ErrorLog

	engine, err := xorm.NewEngine("mysql", conf.Mysql.DataSource)
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{Addr: conf.Redis.DataSource, Password: conf.Redis.Auth})
	rpcxClient, err := grpcx.MustNewGrpcxClient(conf.IntegralRpc)
	if err != nil {
		log.Fatal(err)
	}
	integralRpcModel := integralrpcmodel.NewIntegralRpcModel(
		rpcxClient,
	)

	goodsModel := model.NewGoodsModel(engine, client, conf.Mysql.Table.Goods)
	goodsLogic := logic.NewGoodsLogic(goodsModel, integralRpcModel)
	userController := controller.NewGoodsController(goodsLogic)
	loginAuth := middleware.NewAuthorization(client)
	r := gin.Default()
	r.Use(loginAuth.Auth)
	goodsRouteGroup := r.Group("/goods")
	{
		goodsRouteGroup.POST("/search", userController.GoodSearch)
		goodsRouteGroup.POST("/list", userController.GoodSearch)
		goodsRouteGroup.POST("/order", userController.GoodsOrder)

	}
	log4g.Error(r.Run(conf.Port)) // listen and serve on 0.0.0.0:8080
}
