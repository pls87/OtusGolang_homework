package migrations

import (
	"log"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func Migrate(cfg config.StorageConf) {
	db, err := goose.OpenDBWithDriver("postgres", cfg.ConnString)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	if err := goose.Run("up", db, "."); err != nil {
		log.Fatalf("goose up: %v", err)
	}
}
