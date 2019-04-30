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
	"integral-mall/common/rpcxclient/userrpcmodel"
	"integral-mall/order/command/api/config"
	"integral-mall/order/controller"
	"integral-mall/order/logic"
	"integral-mall/order/model"
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

	rpcxClient, err := grpcx.MustNewGrpcxClient(conf.UserRpc)
	if err != nil {
		log.Fatal(err)
	}
	userRpcModel := userrpcmodel.NewUserRpcModel(
		rpcxClient,
	)
	orderModel := model.NewOrderModel(engine, client, conf.Mysql.Table.Order)
	orderLogic := logic.NewOrderLogic(orderModel, client, userRpcModel)
	orderController := controller.NewOrderController(orderLogic)

	auth := middleware.NewAuthorization(client)
	r := gin.Default()
	userRouteGroup := r.Group("/order")
	userRouteGroup.Use(auth.Auth)
	{
		userRouteGroup.POST("/list", orderController.OrderList)
	}
	log4g.Error(r.Run(conf.Port)) // listen and serve on 0.0.0.0:8080
}
