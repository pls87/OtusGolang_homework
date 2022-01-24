package storage

import (
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/sql"
)

func New(cfg configs.StorageConf) basic.Storage {
	switch cfg.Type {
	case "sql":
		return sql.New(cfg)
	case "memory":
		return memory.New(cfg)
	default:
		return memory.New(cfg)
	}
}
