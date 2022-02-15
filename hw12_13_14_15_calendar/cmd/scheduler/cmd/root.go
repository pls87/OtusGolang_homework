package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/notifications"
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
}

var rc *RootCMD

func (rc *RootCMD) shutDown() {
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

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	rc.logg.Info("connecting to storage...")
	if err := rc.storage.Init(ctx); err != nil {
		rc.logg.Error("failed to connect to storage: " + err.Error())
		cancel()
		os.Exit(1)
	}

	rc.logg.Info("connecting to queue...")
	if err := rc.producer.Init(); err != nil {
		rc.logg.Error("failed to connect to queue: " + err.Error())
		cancel()
		os.Exit(1)
	}

	tick := time.Tick(2 * time.Second)

	defer rc.shutDown()
	for {
		select {
		case <-tick:
			err := rc.removeObsolete()
			if err != nil {
				rc.logg.Errorf("couldn't remove obsolete events: %s", err.Error())
			}
			err = rc.sendNotifications()
			if err != nil {
				rc.logg.Errorf("couldn't send notifications: %s", err.Error())
			}
		case <-ctx.Done():
			return
		}
	}
}

func (rc *RootCMD) init() {
	rc.cfg = configs.New(rc.cfgFile)
	rc.logg = logger.New(rc.cfg.Logger)
}

func (rc *RootCMD) removeObsolete() error {
	return rc.storage.Events().DeleteObsolete(context.Background(), 24*365*time.Hour)
}

func (rc *RootCMD) sendNotifications() error {
	events, err := rc.storage.Events().Select().ToNotify().Execute(context.Background())
	if err != nil {
		rc.logg.Errorf("couldn't get events to notify: %s", err.Error())
		return fmt.Errorf("couldn't get events to notify: %w", err)
	}

	for events.Next() {
		cur, e := events.Current()
		if e != nil {
			rc.logg.Errorf("couldn't get next event: %s", e.Error())
			return fmt.Errorf("couldn't handle notification: %w", e)
		}
		e = rc.producer.Produce(notifications.Message{
			ID:    int(cur.ID),
			Title: cur.Title,
			User:  int(cur.UserID),
			Time:  time.Now(),
		}, false)

		if e != nil {
			rc.logg.Errorf("couldn't send notification to queue: %s", e.Error())
			return fmt.Errorf("couldn't send notification to queue: %w", e)
		}

		e = rc.storage.Notifications().TrackSent(context.Background(), cur.ID)
		if e != nil {
			rc.logg.Errorf("couldn't mark notification sent: %s", e.Error())
			return fmt.Errorf("couldn't mark notification sent: %w", e)
		}
	}
	return nil
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
