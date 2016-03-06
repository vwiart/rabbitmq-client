package config

type Config struct {
	RabbitMQServerURL string
}

func GetConfig() Config {
	return Config{RabbitMQServerURL: "amqp://guest:guest@192.168.0.1:5672/"}
}

