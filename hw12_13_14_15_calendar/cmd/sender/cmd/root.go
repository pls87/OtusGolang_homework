package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/notifications"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const consumerTag = "calendar_sender"

type RootCMD struct {
	*cobra.Command
	cfgFile string
	cfg     configs.Config
	logg    *logrus.Logger

	consumer notifications.Consumer
}

var rc *RootCMD

func (rc *RootCMD) shutDown() {
	if err := rc.consumer.Dispose(); err != nil {
		rc.logg.Errorf("error while consumer shut down: %s", err)
	}
}

func (rc *RootCMD) run() {
	rc.consumer = notifications.NewConsumer(rc.cfg.Queue)
	rc.logg.Info(rc.cfg.Queue)

	ctx, _ := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	rc.logg.Info("connecting to queue...")
	if err := rc.consumer.Init(); err != nil {
		rc.logg.Error("failed to connect to queue: " + err.Error())
		os.Exit(1)
	}

	messages, errors, err := rc.consumer.Consume(consumerTag)
	if err != nil {
		rc.logg.Errorf("couldn't connect to queue: %s", err)
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
			rc.logg.Errorf("error from consumer: %s", e)
		case m, ok = <-messages:
			if !ok {
				break
			}
			rc.logg.Infof("Received message: %v", m)
		case <-ctx.Done():
			ok = false
		}
	}

	rc.shutDown()
}

func (rc *RootCMD) init() {
	rc.cfg = configs.New(rc.cfgFile)
	rc.logg = logger.New(rc.cfg.Logger)
}

func newRootCommand() *RootCMD {
	cmd := new(RootCMD)
	cmd.Command = &cobra.Command{
		Use:   "calendar_scheduler",
		Short: "Scheduler process: generates notifications and removed obsolete events",
		Run: func(c *cobra.Command, args []string) {
			cmd.run()
		},
	}
	return cmd
}

func Execute() error {
	return rc.Execute()
}

func init() {
	cmd := newRootCommand()
	cobra.OnInitialize(cmd.init)
	cmd.PersistentFlags().StringVar(&cmd.cfgFile, "config", "", "config file")
}
