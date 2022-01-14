package sqlstorage

import (
	"context"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/config"

	"github.com/jmoiron/sqlx"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage"

	_ "github.com/lib/pq"
)

type SQLStorage struct {
	cfg    config.StorageConf
	db     *sqlx.DB
	events SQLEventRepository
}

func New(cfg config.StorageConf) *SQLStorage {
	return &SQLStorage{
		events: SQLEventRepository{},
		cfg:    cfg,
	}
}

func (s *SQLStorage) Events() storage.EventRepository {
	return &s.events
}

func (s *SQLStorage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, s.cfg.Driver, s.cfg.ConnString)
	if err == nil {
		s.db = db
		s.events.Attach(s.db)
	}
	return err
}

func (s *SQLStorage) Close() error {
	return s.db.Close()
}
