package common

import (
	"time"

	"github.com/spf13/viper"
)

type loggerConfig struct {
	Level string `mapstructure:"level"`
}

type storageConfig struct {
	Type          string `mapstructure:"type"`
	DSN           string `mapstructure:"dsn"`
	MigrationsDir string `mapstructure:"migrations_dir"`
}

type serverConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type AmqpConfig struct {
	Uri                       string        `mapstructure:"uri"`
	ExchangeName              string        `mapstructure:"exchangeName"`
	ExchangeType              string        `mapstructure:"exchangeType"`
	RoutingKey                string        `mapstructure:"routingKey"`
	QueueName                 string        `mapstructure:"queueName"`
	ConsumerTag               string        `mapstructure:"consumerTag"`
	Reliable                  bool          `mapstructure:"reliable"`
	NotificationPeriodSeconds time.Duration `mapstructure:"notificationPeriodSeconds"`
}

type config struct {
	Logger     loggerConfig  `mapstructure:"logger"`
	Storage    storageConfig `mapstructure:"storage"`
	HttpServer serverConfig  `mapstructure:"httpserver"`
	GrpcServer serverConfig  `mapstructure:"grpcserver"`
	Amqp       AmqpConfig    `mapstructure:"amqp"`
}

var Config = &config{}

func Init(configPath string) error {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(Config)
	if err != nil {
		return err
	}
	return nil
}
