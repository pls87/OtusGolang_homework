package run

import (
	"log"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates data in sql storage",
	Long:  `<Long version desc>`,
	Run: func(cmd *cobra.Command, args []string) {
		if cfg.Storage.Type != "sql" {
			return
		}

		db, err := goose.OpenDBWithDriver(cfg.Storage.Driver, cfg.Storage.ConnString)
		if err != nil {
			log.Fatalf("goose: failed to open DB: %v\n", err)
		}

		defer db.Close()

		if err := goose.Run("up", db, "../migrations"); err != nil {
			log.Printf("goose up: %v", err)
		}
	},
}
