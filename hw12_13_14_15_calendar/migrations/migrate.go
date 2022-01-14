package migrations

import (
	"fmt"
	"log"

	// init postgres driver.
	_ "github.com/lib/pq"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/config"
	"github.com/pressly/goose/v3"
)

func Migrate(cfg config.StorageConf) {
	db, err := goose.OpenDBWithDriver("postgres", cfg.ConnString)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer db.Close()

	if err := goose.Run("up", db, "."); err != nil {
		fmt.Printf("goose up: %v", err)
	}
}
