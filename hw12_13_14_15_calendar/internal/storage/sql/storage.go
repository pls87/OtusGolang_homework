package sql

import (
	"context"

	"github.com/jmoiron/sqlx"

	// init postgres driver.
	_ "github.com/lib/pq"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
)

type Storage struct {
	cfg    configs.StorageConf
	db     *sqlx.DB
	events *EventRepository
}

func New(cfg configs.StorageConf) *Storage {
	return &Storage{
		events: &EventRepository{},
		cfg:    cfg,
	}
}

func (s *Storage) Events() basic.EventRepository {
	return s.events
}

func (s *Storage) Init(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, s.cfg.Driver, s.cfg.Conn)
	if err == nil {
		s.db = db
		s.events.db = s.db
	}
	return err
}

func (s *Storage) Dispose() error {
	return s.db.Close()
}
