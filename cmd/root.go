package cmd

import (
	"log/slog"
	"os"

	"github.com/getsentry/sentry-go"
	slogmulti "github.com/samber/slog-multi"
	slogsentry "github.com/samber/slog-sentry/v2"
	"github.com/spf13/cobra"

	"github.com/nhymxu/go-boilerplate/pkg/config"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sample-project", // TODO: change project name here
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(dependencyInit)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $APPLICATION_DIR/.env)")
	//rootCmd.PersistentFlags().Bool()

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func dependencyInit() {
	err := config.LoadConfig(cfgFile)
	if err != nil {
		panic("Can't load config from environment")
	}

	initLog()
	initSentry()
}

func newBaseHandler() slog.Handler {
	if config.ENV.Debug {
		return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
}

func initLog() {
	slog.SetDefault(slog.New(newBaseHandler()))
}

func initSentry() {
	if config.ENV.Sentry.DSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              config.ENV.Sentry.DSN,
			AttachStacktrace: true,
		})
		if err != nil {
			slog.Error("Sentry initialization failed", "error", err)
		} else {
			slog.Info("Initialized Sentry integration.")
			integrateSlogWithSentry()
		}
	} else {
		slog.Info("SENTRY_DSN not found, sentry integration disabled.")
	}
}

func integrateSlogWithSentry() {
	sentryHandler := slogsentry.Option{
		Level: slog.LevelError,
	}.NewSentryHandler()

	handler := slogmulti.Fanout(newBaseHandler(), sentryHandler)
	slog.SetDefault(slog.New(handler))
	slog.Info("Integrate slog with Sentry successfully.")
}
