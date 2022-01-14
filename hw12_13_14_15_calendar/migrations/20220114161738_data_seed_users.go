package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upDataSeedUsers, downDataSeedUsers)
}

func upDataSeedUsers(tx *sql.Tx) error {
	query := `INSERT INTO users (first_name, last_name, email) VALUES ('Pavel', 'Lysenko', 'plysenko@mail.lo')`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}

func downDataSeedUsers(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
