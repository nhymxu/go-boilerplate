package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/spf13/cobra"

	"github.com/nhymxu/go-boilerplate/apps/api"
	"github.com/nhymxu/go-boilerplate/pkg/config"
)

var apiServerCmd = &cobra.Command{
	Use:   "api",
	Short: "API server",
	Long: `Serve all service on same pod.
Can scale later.`,
	Run: func(cmd *cobra.Command, _ []string) {
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			host = ""
		}

		port, err := cmd.Flags().GetInt64("port")
		if err != nil {
			panic("Get port config error")
		}

		shutdownTime, err := cmd.Flags().GetInt64("shutdown_time")
		if err != nil {
			shutdownTime = config.APIDefaultShutdownTime
		}

		e := api.New()

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		sc := echo.StartConfig{
			Address:         fmt.Sprintf("%s:%d", host, port),
			HideBanner:      true,
			GracefulTimeout: time.Duration(shutdownTime) * time.Second,
			OnShutdownError: func(err error) {
				slog.Error("graceful shutdown error", "error", err)
			},
		}

		if err := sc.Start(ctx, e); err != nil {
			slog.Error("shutting down the server", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apiServerCmd)

	apiServerCmd.Flags().String("host", "", "API host listening")
	apiServerCmd.Flags().Int64("port", 8000, "API port listening")
	apiServerCmd.Flags().Int64("shutdown_time", config.APIDefaultShutdownTime, "Gracefully shutdown time")
}
