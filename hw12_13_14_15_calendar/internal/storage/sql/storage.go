package sqlstorage

import (
	"context"

	"github.com/jmoiron/sqlx"

	// init postgres driver.
	_ "github.com/lib/pq"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	basicstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
)

type SQLStorage struct {
	cfg    configs.StorageConf
	db     *sqlx.DB
	events SQLEventRepository
}

func New(cfg configs.StorageConf) *SQLStorage {
	return &SQLStorage{
		events: SQLEventRepository{},
		cfg:    cfg,
	}
}

func (s *SQLStorage) Events() basicstorage.EventRepository {
	return &s.events
}

func (s *SQLStorage) Init(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, s.cfg.Driver, s.cfg.Conn)
	if err == nil {
		s.db = db
		s.events.Attach(s.db)
	}
	return err
}

func (s *SQLStorage) Dispose() error {
	return s.db.Close()
}
