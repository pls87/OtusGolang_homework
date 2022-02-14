package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/cmd/shared"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/notifications"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const consumerTag = "calendar_sender"

var (
	cfg     configs.Config
	cfgFile string
	logg    *logrus.Logger

	rootCmd = &cobra.Command{
		Use:   "calendar_scheduler",
		Short: "A background process to send notifications",
		Run: func(cmd *cobra.Command, args []string) {
			consumer := notifications.NewConsumer(cfg.Queue)
			logg.Info(cfg.Queue)

			ctx, cancel := signal.NotifyContext(context.Background(),
				syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
			defer cancel()

			shutDown := func() {
				if err := consumer.Dispose(); err != nil {
					logg.Errorf("error while consumer shut down: %s", err)
				}
			}

			logg.Info("connecting to queue...")
			if err := consumer.Init(); err != nil {
				logg.Error("failed to connect to queue: " + err.Error())
				cancel()
				os.Exit(1)
			}

			messages, errors, err := consumer.Consume(consumerTag)
			if err != nil {
				logg.Errorf("couldn't connect to queue: %s", err)
				cancel()
				os.Exit(1)
			}

			var e error
			var m notifications.Message
			for ok := true; ok; {
				select {
				case e, ok = <-errors:
					if !ok {
						break
					}
					logg.Errorf("error from consumer: %s", e)
				case m, ok = <-messages:
					if !ok {
						break
					}
					logg.Infof("Received message: %v", m)
				case <-ctx.Done():
					ok = false
				}
			}

			shutDown()
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
