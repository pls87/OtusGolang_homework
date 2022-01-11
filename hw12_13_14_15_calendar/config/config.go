package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf
}

type LoggerConf struct {
	Level string
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
	viper.SetDefault("loglevel", "debug")

	cfg.Logger.Level = viper.GetString("loglevel")

	return cfg
}
