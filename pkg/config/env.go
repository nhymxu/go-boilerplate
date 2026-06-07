package config

import (
	"github.com/jinzhu/copier"
	"github.com/nhymxu/gommon/cfgloader"
)

// Config define mapping struct field and environment field
type Config struct {
	Debug  bool `koanf:"DEBUG"`
	Sentry struct {
		DSN string `koanf:"DSN"`
	} `koanf:"SENTRY"`

	TokenAuth string `koanf:"TOKEN_AUTH" copier:"-"`
}

func (c *Config) Sanitized() Config {
	var cc Config

	// Secrets excluded ❌
	err := copier.Copy(&cc, &c)
	if err != nil {
		return Config{}
	}

	return cc
}

// C is global variable for using config in other place
var C Config

// Load read env file and loaded to environment and global ENV variable
func Load(cfgFile string) error {
	var err error

	C, err = cfgloader.LoadConfig[Config](cfgFile, configDefaults)

	return err
}
