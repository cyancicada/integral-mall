package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"integral-mall/common/utils"
	"integral-mall/order/command/rpc/config"
	"integral-mall/order/logic"
	"integral-mall/order/model"
	"integral-mall/order/protos"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/yakaa/grpcx"
	"github.com/yakaa/log4g"
	"google.golang.org/grpc"
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
	engine, err := xorm.NewEngine("mysql", conf.Mysql.DataSource)
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{Addr: conf.Redis.DataSource, Password: conf.Redis.Auth})
	orderModel := model.NewOrderModel(engine, client, conf.Mysql.Table.Order)

	orderQunue, err := utils.NewRabbitMqServer(
		conf.RabbitMq.DataSource+conf.RabbitMq.VirtualHost,
		conf.RabbitMq.QueueName,
	)
	if err != nil {
		log.Fatal(err)
	}

	orderServerLogic := logic.NewOrderRpcServerLogic(orderModel, orderQunue)
	rpcServer, err := grpcx.MustNewGrpcxServer(conf.RpcServerConfig, func(server *grpc.Server) {
		protos.RegisterOrderRpcServer(server, orderServerLogic)
	})
	if err != nil {
		log.Fatal(err)
	}
	orderServerLogic.ConsumeMessageStart()
	defer orderServerLogic.CloseRabbitMqServer()
	log4g.Error(rpcServer.Run())
}
