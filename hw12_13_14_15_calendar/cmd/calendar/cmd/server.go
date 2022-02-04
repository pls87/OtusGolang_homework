package calendarcmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts server. gRPC/HTTP and host/port are in the config",
	Run: func(cmd *cobra.Command, args []string) {
		storage := storage.New(cfg.Storage)
		calendar := app.New(logg, storage)

		server := server.New(logg, calendar, cfg.API)

		ctx, cancel := signal.NotifyContext(context.Background(),
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		go func() {
			<-ctx.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			if err := storage.Dispose(); err != nil {
				logg.Error("failed to close storage connection: " + err.Error())
			}

			if err := server.Stop(ctx); err != nil {
				logg.Error("failed to stop http internal: " + err.Error())
			}
		}()

		logg.Info("connecting to storage...")

		if err := storage.Init(ctx); err != nil {
			logg.Error("failed to connect to storage: " + err.Error())
			cancel()
			os.Exit(1)
		}

		logg.Info("calendar is running...")

		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start server: " + err.Error())
			cancel()
			os.Exit(1)
		}

		<-ctx.Done()
	},
}
