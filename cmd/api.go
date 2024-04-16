package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"

	"github.com/nhymxu/go-boilerplate/apps/api"
)

var apiServerCmd = &cobra.Command{
	Use:   "api",
	Short: "API server",
	Long: `Serve all service on same pod.
Can scale later.`,
	Run: func(cmd *cobra.Command, _ []string) {
		port, err := cmd.Flags().GetInt64("port")
		if err != nil {
			panic("Get port config error")
		}

		e := api.New()
		if err != nil {
			panic("Something wrong")
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()
		// Start server
		go func() {
			if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
				e.Logger.Fatalf("shutting down the server. Err: %v", err)
			}
		}()

		// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(apiServerCmd)

	apiServerCmd.Flags().Int64("port", 8000, "API port listening")
}
