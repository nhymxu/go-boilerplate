package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// configDefaults holds default values for all config fields.
// Nested struct fields use dot-separated koanf tag paths (e.g. "SENTRY.DSN").
var configDefaults = map[string]any{
	"DEBUG":      false,
	"SENTRY.DSN": "",
	"TOKEN_AUTH": "",
}

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

	// Load defaults first (lowest precedence)
	if err := k.Load(confmap.Provider(configDefaults, "."), nil); err != nil {
		return fmt.Errorf("failed to load default config: %w", err)
	}

	configFile := ".env"
	if cfgFile != "" {
		configFile = cfgFile
	}

	// Load from config file (optional – skip if file does not exist)
	if err := k.Load(file.Provider(configFile), dotenv.Parser()); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("failed to load config file %s: %w", configFile, err)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Using config file:", configFile)
	}

	// Override with actual environment variables (highest precedence)
	if err := k.Load(env.Provider("", ".", nil), nil); err != nil {
		return err
	}

	return k.Unmarshal("", &ENV)
}
