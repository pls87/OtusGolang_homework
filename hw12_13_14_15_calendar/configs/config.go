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
)

type Config struct {
	Logger  LoggerConf  `toml:"logger"`
	Storage StorageConf `toml:"storage"`
	API     APIConf     `toml:"api"`
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

func New(cfgFile string) Config {
	cfg := Config{
		Logger:  LoggerConf{Level: DefaultLogLevel},
		Storage: StorageConf{Type: DefaultStorageType},
		API:     APIConf{Type: DefaultAPIType, Port: DefaultAPIPort},
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
