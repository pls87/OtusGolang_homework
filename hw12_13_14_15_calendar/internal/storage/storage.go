package storage

import (
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	basicstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	memorystorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/sql"
)

func New(cfg configs.StorageConf) basicstorage.Storage {
	switch cfg.Type {
	case "sql":
		return sqlstorage.New(cfg)
	case "memory":
		return memorystorage.New(cfg)
	default:
		return memorystorage.New(cfg)
	}
}
