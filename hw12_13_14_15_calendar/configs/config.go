package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	DefaultLogLevel    = "debug"
	DefaultStorageType = "memory"
	DefaultAPIType     = "http"
	DefaultAPIPort     = 8082
	DefaultQueuePort   = 5672
	DefaultQueueHost   = "127.0.0.1"
)

type Config struct {
	Logger       LoggerConf       `mapstructure:"logger"`
	Storage      StorageConf      `mapstructure:"storage"`
	API          APIConf          `mapstructure:"api"`
	Notification NotificationConf `mapstructure:"notifications"`
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
}

type StorageConf struct {
	Type   string `mapstructure:"type"`
	Driver string `mapstructure:"driver"`
	Conn   string `mapstructure:"conn"`
}

type APIConf struct {
	Type string `mapstructure:"type"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type NotificationConf struct {
	User     string `mapstructure:"rabbit_user"`
	Password string `mapstructure:"rabbit_password"`
	Host     string `mapstructure:"rabbit_host"`
	Port     int    `mapstructure:"rabbit_port"`
	Interval int    `mapstructure:"scheduler_interval"`
}

func New(cfgFile string) Config {
	cfg := Config{
		Logger:  LoggerConf{Level: DefaultLogLevel},
		Storage: StorageConf{Type: DefaultStorageType},
		API:     APIConf{Type: DefaultAPIType, Port: DefaultAPIPort},
		Notification: NotificationConf{
			Port: DefaultQueuePort,
			Host: DefaultQueueHost,
		},
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			viper.Unmarshal(&cfg)
			fmt.Println("Using configs file:", viper.ConfigFileUsed())
		}
	}

	return cfg
}
