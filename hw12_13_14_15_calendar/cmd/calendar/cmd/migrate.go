package calendarcmd

import (
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/migrations"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates data in sql storage",
	Run: func(cmd *cobra.Command, args []string) {
		migrations.Migrate(cfg.Storage)
	},
}
