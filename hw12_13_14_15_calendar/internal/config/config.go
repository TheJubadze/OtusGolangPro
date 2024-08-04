package config

import (
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

type config struct {
	Logger     loggerConfig  `mapstructure:"logger"`
	Storage    storageConfig `mapstructure:"storage"`
	HttpServer serverConfig  `mapstructure:"httpserver"`
	GrpcServer serverConfig  `mapstructure:"grpcserver"`
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
