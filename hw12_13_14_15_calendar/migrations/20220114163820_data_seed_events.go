package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upDataSeedEvents, downDataSeedEvents)
}

func upDataSeedEvents(tx *sql.Tx) error {
	query := `INSERT INTO "events" (title, user_id, start, duration, notify_before,  description) VALUES 
		('Lunch', 1, TIMESTAMP WITH TIME ZONE '2022-02-15 7:00:00+00', '1 hour', '30 minutes', 'Time to eat'),
		('Daily Scrum', 1, TIMESTAMP WITH TIME ZONE '2022-02-15 10:00:00+00', '30 minutes', '30 minutes', 'Time to meet')`
	_, err := tx.Exec(query)
	return err
}

func downDataSeedEvents(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
