package calendarcmd

import (
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/migrations"
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
		migrations.Migrate(cfg.Storage)
	},
}
