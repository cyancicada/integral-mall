package config

import (
	rpc_config "github.com/yakaa/grpcx/config"
)

type Config struct {
	RpcServerConfig *rpc_config.ServiceConf
	Mode            string `json:"mode"`
	Port            string `json:"port"`
	Mysql           struct {
		DataSource string
		Table      struct {
			Integral string
		}
	}
	Redis struct {
		DataSource string
		Auth       string
	}
	RabbitMq struct {
		DataSource  string
		VirtualHost string
		QueueName   string
	}
}
