package app

import (
	"github.com/urfave/cli/v2"
)

const (
	appName            = "cryptoexchange"
	defaultLoggerLevel = "info"
)

type Config struct {
	LogLevel string
	HTTP     *WebConfig
}

func NewConfig() *Config {
	return &Config{
		HTTP: NewWebConfig(),
	}
}

// BuildFlags bindings.
func (config *Config) BuildFlags() []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			EnvVars:     []string{"LOG_LEVEL"},
			Value:       defaultLoggerLevel,
			Destination: &config.LogLevel,
		},
	}

	flags = append(flags, config.HTTP.BuildFlags()...)

	return flags
}
