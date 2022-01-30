package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upDataSeedEvents, downDataSeedEvents)
}

func upDataSeedEvents(tx *sql.Tx) error {
	query := `INSERT INTO "events" (title, user_id, start, duration, notify_before,  description) 
		VALUES ('Lunch', 1, TIMESTAMP WITH TIME ZONE '2022-01-15 13:00:00+06', '1 hour', '30 minutes', 'Eat and eat again')`
	_, err := tx.Exec(query)
	return err
}

func downDataSeedEvents(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
