package cmd

import (
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/cmd/shared"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cfg     configs.Config
	cfgFile string
	logg    *logrus.Logger

	rootCmd = &cobra.Command{
		Use:   "calendar_scheduler",
		Short: "A background process to generate notifications",
		Run: func(cmd *cobra.Command, args []string) {
			logg.Info("Root command does nothing. Please use special commands")
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(beforeRun)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.AddCommand(shared.VersionCmd)
}

func beforeRun() {
	cfg = configs.New(cfgFile)
	logg = logger.New(cfg.Logger)
}
