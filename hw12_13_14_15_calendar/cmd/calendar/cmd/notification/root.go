package notification

import (
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cfg  configs.Config
	logg *logrus.Logger
)

var Cmd = &cobra.Command{
	Use:   "notification",
	Short: "Root command for ",
	Run: func(cmd *cobra.Command, args []string) {
		logg.Info("Root notification command does nothing. Please use generate/send subcommands")
	},
}

func init() {
	Cmd.AddCommand(generateCmd, sendCmd)
}

func SetConfig(c configs.Config, log *logrus.Logger) {
	cfg = c
	logg = log
}
