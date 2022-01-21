package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf  `toml:"logger"`
	Storage StorageConf `toml:"storage"`
	Net     NetConf     `toml:"net"`
}

type LoggerConf struct {
	Level string `toml:"level"`
}

type StorageConf struct {
	Type   string `toml:"type"`
	Driver string `toml:"driver"`
	Conn   string `toml:"conn"`
}

type NetConf struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func New(cfgFile string) Config {
	cfg := Config{
		Logger:  LoggerConf{Level: "debug"},
		Storage: StorageConf{Type: "memory"},
		Net:     NetConf{Port: 8082},
	}
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using configs file:", viper.ConfigFileUsed())
	}

	viper.AutomaticEnv()
	viper.Unmarshal(&cfg)

	return cfg
}
