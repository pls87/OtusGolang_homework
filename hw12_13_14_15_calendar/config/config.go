package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	Type       string
	Driver     string
	ConnString string
}

func New(cfgFile string) Config {
	cfg := Config{}
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.AutomaticEnv()
	viper.SetDefault("log.level", "debug")
	viper.SetDefault("storage.type", "memory")

	cfg.Logger.Level = viper.GetString("log.level")
	cfg.Storage.Type = viper.GetString("storage.type")
	cfg.Storage.Driver = viper.GetString("storage.driver")
	cfg.Storage.ConnString = viper.GetString("storage.conn")

	return cfg
}
