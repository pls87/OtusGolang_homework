package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upNotificationSent, downNotificationSent)
}

func upNotificationSent(tx *sql.Tx) error {
	query := `CREATE TABLE "notification_sent"
	(
    	"ID"            serial         NOT NULL,
    	"event_id"      integer 	   NOT NULL UNIQUE,
    	"time" 			timestamptz    NOT NULL,
    	CONSTRAINT "notification_sent_ID" PRIMARY KEY ("ID")
	);
	ALTER TABLE ONLY "notification_sent"
    	ADD CONSTRAINT "notification_sent_event_id_fkey" FOREIGN KEY (event_id) REFERENCES "events" ("ID") ON DELETE CASCADE;`

	_, err := tx.Exec(query)
	return err
}

func downNotificationSent(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE "notification_sent"`)
	return err
}
