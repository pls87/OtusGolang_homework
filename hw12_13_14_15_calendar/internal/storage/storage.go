package storage

import (
	"context"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/config"
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
