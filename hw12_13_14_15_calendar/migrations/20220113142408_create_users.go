package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateUsers, downCreateUsers)
}

func upCreateUsers(tx *sql.Tx) error {
	query := `CREATE TABLE "users"(
    "ID"         serial         NOT NULL,
    "first_name" character(255) NOT NULL,
    "last_name"  character(255) NOT NULL,
    "email"      character(255) NOT NULL,
    CONSTRAINT "users_ID" PRIMARY KEY ("ID")
)`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}

func downCreateUsers(tx *sql.Tx) error {
	if _, err := tx.Exec(`DROP TABLE "users"`); err != nil {
		return err
	}
	return nil
}
