package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/notifications"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/sender"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type RootCMD struct {
	*cobra.Command
	cfgFile string
	cfg     configs.Config
	logg    *logrus.Logger

	consumer notifications.Consumer

	sender sender.Sender
}

var rc *RootCMD

func (rc *RootCMD) shutDown() {
	rc.sender.Stop()
	if err := rc.consumer.Dispose(); err != nil {
		rc.logg.Errorf("error while consumer shut down: %s", err)
	}
}

func (rc *RootCMD) messageHandler(m notifications.Message) {
	rc.logg.Infof("received message: %v", m)
}

func (rc *RootCMD) errorHandler(e error) {
	rc.logg.Errorf("error while consuming message: %s", e)
}

func (rc *RootCMD) run() {
	rc.consumer = notifications.NewConsumer(rc.cfg.Queue)

	ctx, _ := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	rc.logg.Info("connecting to queue...")
	if err := rc.consumer.Init(); err != nil {
		rc.logg.Error("failed to connect to queue: " + err.Error())
		os.Exit(1)
	}

	rc.sender = sender.NewSender(rc.consumer, rc.messageHandler, rc.errorHandler)
	if err := rc.sender.Send(); err != nil {
		rc.logg.Errorf("couldn't consume messages: %s", err)
		os.Exit(1)
	}
	<-ctx.Done()
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
