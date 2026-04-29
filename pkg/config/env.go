package config

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// EnvConfigMap define mapping struct field and environment field
type EnvConfigMap struct {
	Debug  bool `koanf:"DEBUG"`
	Sentry struct {
		DSN string `koanf:"DSN"`
	} `koanf:"SENTRY"`

	TokenAuth string `koanf:"TOKEN_AUTH"`
}

// ENV is global variable for using config in other place
var ENV EnvConfigMap

// LoadConfig read env file and loaded to environment and global ENV variable
func LoadConfig(cfgFile string) error {
	k := koanf.New(".")

	configFile := ".env"
	if cfgFile != "" {
		configFile = cfgFile
	}

	// Load from config file
	if err := k.Load(file.Provider(configFile), dotenv.Parser()); err != nil {
		return fmt.Errorf("failed to load config file %s: %w", configFile, err)
	}
	fmt.Fprintln(os.Stderr, "Using config file:", configFile)

	// Override with actual environment variables
	if err := k.Load(env.Provider("", ".", nil), nil); err != nil {
		return err
	}

	return k.Unmarshal("", &ENV)
}
