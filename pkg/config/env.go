package config

import "github.com/nhymxu/gommon/cfgloader"

// EnvConfigMap define mapping struct field and environment field
type EnvConfigMap struct {
	Debug  bool `koanf:"DEBUG"`
	Sentry struct {
		DSN string `koanf:"DSN"`
	} `koanf:"SENTRY"`

	TokenAuth string `koanf:"TOKEN_AUTH"`
}

// C is global variable for using config in other place
var C EnvConfigMap

// LoadConfig read env file and loaded to environment and global ENV variable
func LoadConfig(cfgFile string) error {
	var err error
	C, err = cfgloader.LoadConfig[EnvConfigMap](cfgFile, configDefaults)
	return err
}
