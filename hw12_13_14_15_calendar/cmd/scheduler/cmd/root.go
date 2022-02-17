package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/notifications"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type RootCMD struct {
	*cobra.Command
	cfgFile string
	cfg     configs.Config
	logg    *logrus.Logger

	storage  basic.Storage
	producer notifications.Producer

	scheduler scheduler.Scheduler
}

var rc *RootCMD

func (rc *RootCMD) shutDown() {
	rc.scheduler.Stop()

	if err := rc.storage.Dispose(); err != nil {
		rc.logg.Error("failed to close storage connection: " + err.Error())
	}
	if err := rc.producer.Dispose(); err != nil {
		rc.logg.Error("failed to close producer connection: " + err.Error())
	}
}

func (rc *RootCMD) run() {
	rc.storage = storage.New(rc.cfg.Storage)
	rc.producer = notifications.NewProducer(rc.cfg.Queue)

	rc.scheduler = scheduler.NewScheduler(rc.producer, rc.storage)

	ctx, _ := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	rc.logg.Info("connecting to storage...")
	if err := rc.storage.Init(ctx); err != nil {
		rc.logg.Error("failed to connect to storage: " + err.Error())
		os.Exit(1)
	}

	rc.logg.Info("connecting to queue...")
	if err := rc.producer.Init(); err != nil {
		rc.logg.Error("failed to connect to queue: " + err.Error())
		os.Exit(1)
	}

	errors := rc.scheduler.Start(2 * time.Second)
	defer rc.shutDown()
	for {
		select {
		case e := <-errors:
			rc.logg.Errorf("error while handling event: %s", e)
		case <-ctx.Done():
			return
		}
	}
}

func (rc *RootCMD) init() {
	rc.cfg = configs.New(rc.cfgFile)
	rc.logg = logger.New(rc.cfg.Logger)
}

func newRootCommand() *RootCMD {
	cmd := new(RootCMD)
	cmd.Command = &cobra.Command{
		Use:   "calendar_scheduler",
		Short: "Scheduler process: generates notifications and removes obsolete events",
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
