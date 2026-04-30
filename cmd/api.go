package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/urfave/cli/v2"

	"github.com/nhymxu/go-boilerplate/apps/api"
	"github.com/nhymxu/go-boilerplate/pkg/config"
)

func apiCommand() *cli.Command {
	return &cli.Command{
		Name:  "api",
		Usage: "API server",
		Description: `Serve all service on same pod.
Can scale later.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "host",
				Value: "",
				Usage: "API host listening",
			},
			&cli.Int64Flag{
				Name:  "port",
				Value: 8000,
				Usage: "API port listening",
			},
			&cli.Int64Flag{
				Name:  "shutdown_time",
				Value: config.APIDefaultShutdownTime,
				Usage: "Gracefully shutdown time",
			},
		},
		Action: func(c *cli.Context) error {
			host := c.String("host")
			port := c.Int64("port")
			shutdownTime := c.Int64("shutdown_time")

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

			return nil
		},
	}
}
