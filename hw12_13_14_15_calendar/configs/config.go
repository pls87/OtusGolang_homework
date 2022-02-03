package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	DefaultLogLevel        = "debug"
	DefaultStorageType     = "memory"
	DefaultAPIType         = "http"
	DefaultAPIPort         = 8082
	DefaultQueuePort       = 5672
	DefaultQueueExchange   = "calendar"
	DefaultQueueName       = "notification"
	DefaultQueueRoutingKey = "before_event"
)

type Config struct {
	Logger  LoggerConf  `toml:"logger"`
	Storage StorageConf `toml:"storage"`
	API     APIConf     `toml:"api"`
	Queue   QueueConf   `toml:"queue"`
}

type LoggerConf struct {
	Level string `toml:"level"`
}

type StorageConf struct {
	Type   string `toml:"type"`
	Driver string `toml:"driver"`
	Conn   string `toml:"conn"`
}

type APIConf struct {
	Type string `toml:"type"`
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type QueueConf struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Exchange string `toml:"exchange"`
	Queue    string `toml:"queue"`
	Key      string `toml:"key"`
}

func New(cfgFile string) Config {
	cfg := Config{
		Logger:  LoggerConf{Level: DefaultLogLevel},
		Storage: StorageConf{Type: DefaultStorageType},
		API:     APIConf{Type: DefaultAPIType, Port: DefaultAPIPort},
		Queue: QueueConf{
			Port:     DefaultQueuePort,
			Exchange: DefaultQueueExchange,
			Queue:    DefaultQueueName,
			Key:      DefaultQueueRoutingKey,
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
