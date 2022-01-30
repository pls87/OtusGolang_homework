package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	DefaultLogLevel    = "debug"
	DefaultStorageType = "memory"
	DefaultAPIPort     = 8082
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
		Logger:  LoggerConf{Level: DefaultLogLevel},
		Storage: StorageConf{Type: DefaultStorageType},
		Net:     NetConf{Port: DefaultAPIPort},
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
