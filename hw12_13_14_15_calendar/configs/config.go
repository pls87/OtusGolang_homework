package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Net     HTTPConf
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	Type       string
	Driver     string
	ConnString string
}

type HTTPConf struct {
	Host string
	Port int
}

func New(cfgFile string) Config {
	cfg := Config{}
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using configs file:", viper.ConfigFileUsed())
	}

	viper.AutomaticEnv()
	viper.SetDefault("log.level", "debug")
	viper.SetDefault("http.host", "")
	viper.SetDefault("http.port", 8080)
	viper.SetDefault("storage.type", "memory")

	cfg.Logger.Level = viper.GetString("log.level")
	cfg.Storage.Type = viper.GetString("storage.type")
	cfg.Storage.Driver = viper.GetString("storage.driver")
	cfg.Storage.ConnString = viper.GetString("storage.conn")

	cfg.Net.Host = viper.GetString("net.host")
	cfg.Net.Port = viper.GetInt("net.port")

	return cfg
}
