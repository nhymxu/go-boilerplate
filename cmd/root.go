package cmd

import (
	"fmt"
	"os"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var cfgFile string
var l *zap.SugaredLogger

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
	initConfig()
	initLog()
	initSentry()
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		//viper.SetConfigType("yaml")
		//viper.SetConfigName(".obm-bot-crawler")
		viper.SetConfigFile(".env")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func initLog() {
	var logger *zap.Logger
	if viper.GetBool("DEBUG") {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	defer logger.Sync() //nolint:errcheck
	zap.ReplaceGlobals(logger)

	l = zap.S()
}

func initSentry() {
	if sentryDNS := viper.GetString("SENTRY_DSN"); sentryDNS != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDNS,
			AttachStacktrace: true,
			Debug:            true,
		})
		if err != nil {
			l.Errorf("Sentry initialization failed: %v", err)
		} else {
			l.Infof("Initialized Sentry integration.")
			integrateZapWithSentry()
		}
	} else {
		l.Infof("SENTRY_DSN not found, sentry integration disabled.")
	}
}

func integrateZapWithSentry() {
	log := zap.L()
	sentryClient := sentry.CurrentHub().Client()

	cfg := zapsentry.Configuration{
		Level:             zapcore.ErrorLevel, //when to send message to sentry
		EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
		BreadcrumbLevel:   zapcore.InfoLevel,  // at what level should we sent breadcrumbs to sentry
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(sentryClient))

	//in case of err it will return noop core. so we can safely attach it
	if err != nil {
		log.Warn("failed to init zap", zap.Error(err))
	}

	log = zapsentry.AttachCoreToLogger(core, log)

	// to use breadcrumbs feature - create new scope explicitly
	// and attach after attaching the core
	zap.ReplaceGlobals(log.With(zapsentry.NewScope()))

	l = zap.S()
	l.Info("Integrate uber/zap with Sentry successfully.")
}
