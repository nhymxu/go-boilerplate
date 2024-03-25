package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// EnvConfigMap define mapping struct field and environment field
type EnvConfigMap struct {
	Debug  bool `mapstructure:"DEBUG"`
	Sentry struct {
		DSN string `mapstructure:"DSN"`
	}
}

// ENV is global variable for using config in other place
var ENV EnvConfigMap

// LoadConfig read env file and loaded to environment and global ENV variable
func LoadConfig(cfgFile string) error {
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
	err := viper.ReadInConfig()
	if err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		return err
	}

	err = viper.Unmarshal(&ENV)
	if err != nil {
		return err
	}

	return nil
}
