package run

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/migrations"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/config"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/http"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/cobra"
)

var (
	cfg     config.Config
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "calendar",
		Short: "A simple app to manage your events",
		Long:  `<Some long desc here...>`,
		Run: func(cmd *cobra.Command, args []string) {
			logg := logger.New(cfg.Logger.Level)

			storage := newStorage(cfg.Storage)
			calendar := app.New(logg, storage, cfg)

			server := internalhttp.NewServer(logg, calendar)

			ctx, cancel := signal.NotifyContext(context.Background(),
				syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
			defer cancel()

			go func() {
				<-ctx.Done()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				if err := storage.Close(); err != nil {
					logg.Error("failed to close storage connection: " + err.Error())
				}

				if err := server.Stop(ctx); err != nil {
					logg.Error("failed to stop http server: " + err.Error())
				}
			}()

			logg.Info("connecting to storage...")

			if err := storage.Connect(ctx); err != nil {
				logg.Error("failed to connect to storage: " + err.Error())
				cancel()
				os.Exit(1) //nolint:gocritic
			}

			logg.Info("calendar is running...")

			if err := server.Start(ctx); err != nil {
				logg.Error("failed to start http server: " + err.Error())
				cancel()
				os.Exit(1) //nolint:gocritic
			}
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
}

func beforeRun() {
	cfg = config.New(cfgFile)

	if cfg.Storage.Type == "sql" {
		migrations.Migrate(cfg.Storage)
	}
}

func newStorage(cfg config.StorageConf) storage.Storage {
	switch cfg.Type {
	case "sql":
		return sqlstorage.New(cfg)
	case "memory":
		return memorystorage.New(cfg)
	default:
		return memorystorage.New(cfg)
	}
}
