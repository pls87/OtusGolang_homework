package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateEvents, downCreateEvents)
}

func upCreateEvents(tx *sql.Tx) error {
	query := `CREATE TABLE "events"
	(
    	"ID"            serial         NOT NULL,
    	"title"         character(255) NOT NULL,
    	"start"         timestamptz    NOT NULL,
    	"duration"      interval       NOT NULL,
    	"notify_before" interval       NOT NULL,
    	"desc"          text           NOT NULL,
    	"user_id"       integer        NOT NULL,
    	CONSTRAINT "events_ID" PRIMARY KEY ("ID")
	);
	ALTER TABLE ONLY "events"
    	ADD CONSTRAINT "events_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users ("ID") ON DELETE CASCADE;`

	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}

func downCreateEvents(tx *sql.Tx) error {
	if _, err := tx.Exec(`DROP TABLE "events"`); err != nil {
		return err
	}
	return nil
}
