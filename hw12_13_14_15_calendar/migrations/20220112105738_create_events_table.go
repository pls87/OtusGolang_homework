package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateEventsTable, downCreateEventsTable)
}

func upCreateEventsTable(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func downCreateEventsTable(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
