package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/run/config"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage/memory"

	"github.com/spf13/cobra"
)

var (
	cfg config.Config

	rootCmd = &cobra.Command{
		Use:   "calendar",
		Short: "A simple app to manage your events",
		Long:  `<Some long desc here...>`,
		Run: func(cmd *cobra.Command, args []string) {
			logg := logger.New(cfg.Logger.Level)

			storage := memorystorage.New()
			calendar := app.New(logg, storage)

			server := internalhttp.NewServer(logg, calendar)

			ctx, cancel := signal.NotifyContext(context.Background(),
				syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
			defer cancel()

			go func() {
				<-ctx.Done()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				if err := server.Stop(ctx); err != nil {
					logg.Error("failed to stop http server: " + err.Error())
				}
			}()

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
	var cfgFile string
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	cfg = config.Init(cfgFile)
}
