package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/cmd/shared"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/notifications"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/basic"
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
			storage := storage.New(cfg.Storage)
			producer := notifications.NewProducer(cfg.Queue)

			ctx, cancel := signal.NotifyContext(context.Background(),
				syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
			defer cancel()

			shutDown := func() {
				if err := storage.Dispose(); err != nil {
					logg.Error("failed to close storage connection: " + err.Error())
				}
				if err := producer.Dispose(); err != nil {
					logg.Error("failed to close producer connection: " + err.Error())
				}
			}

			logg.Info("connecting to storage...")
			if err := storage.Init(ctx); err != nil {
				logg.Error("failed to connect to storage: " + err.Error())
				cancel()
				os.Exit(1)
			}

			logg.Info("connecting to queue...")
			if err := producer.Init(); err != nil {
				logg.Error("failed to connect to queue: " + err.Error())
				cancel()
				os.Exit(1)
			}

			tick := time.Tick(2 * time.Second)

			defer shutDown()
			for {
				select {
				case <-tick:
					sendNotifications(storage, producer)
				case <-ctx.Done():
					return
				}
			}
		},
	}
)

func sendNotifications(storage basic.Storage, producer notifications.Producer) {
	events, err := storage.Events().Select().ToNotify().Execute(context.Background())
	if err != nil {
		logg.Errorf("couldn't get events to notify: %s", err.Error())
	}

	for events.Next() {
		cur, e := events.Current()
		if e != nil {
			logg.Errorf("couldn't get next event: %s", e.Error())
			return
		}
		e = producer.Produce(notifications.Message{
			ID:    int(cur.ID),
			Title: cur.Title,
			User:  int(cur.UserID),
			Time:  time.Now(),
		}, false)

		if e != nil {
			logg.Errorf("couldn't send notification to queue: %s", e.Error())
			return
		}

		e = storage.Notifications().TrackSent(context.Background(), cur.ID)
		if e != nil {
			logg.Errorf("couldn't mark notification as sent: %s", e.Error())
			return
		}
	}
}

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
