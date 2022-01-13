package storage

import (
	"context"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/config"
	memorystorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/migrations"
)

type ID uint64

type Storage interface {
	Events() EventRepository
	Connect(ctx context.Context) error
	Close() error
}

func Init(cfg config.StorageConf) {
	if cfg.Type == "sql" {
		migrations.Migrate(cfg)
	}
}

func New(cfg config.StorageConf) Storage {
	switch cfg.Type {
	case "sql":
		return sqlstorage.New(cfg)
	case "memory":
		return memorystorage.New(cfg)
	default:
		return memorystorage.New(cfg)
	}
}
