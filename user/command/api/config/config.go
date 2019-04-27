package config

type Config struct {
	Mode  string `json:"mode"`
	Port  string `json:"port"`
	Mysql struct {
		DataSource string
		Table      struct {
			User string
		}
	}
	Redis struct {
		DataSource string
		Auth       string
	}

	RabbitMq string
}
