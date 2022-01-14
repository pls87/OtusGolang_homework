package migrations

import (
	"log"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/config"
	"github.com/pressly/goose/v3"
)

func Migrate(cfg config.StorageConf) {
	if cfg.Type != "sql" {
		return
	}

	db, err := goose.OpenDBWithDriver(cfg.Driver, cfg.ConnString)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer db.Close()

	if err := goose.Run("up", db, "."); err != nil {
		log.Printf("goose up: %v", err)
	}
}
