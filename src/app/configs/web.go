package app

import (
	"time"

	"github.com/urfave/cli/v2"
)

const (
	defaultHTTPPort               = ":3000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultShutdownDuration       = 5 * time.Second
	defaultShutdownTimeout        = 15 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
)

type WebConfig struct {
	Port               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxHeaderMegabytes int
	ShutdownDelay      time.Duration
	ShutdownTimeout    time.Duration
}

func NewWebConfig() *WebConfig {
	return &WebConfig{}
}

// BuildFlags bindings.
func (config *WebConfig) BuildFlags() []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "http-port",
			EnvVars:     []string{"HTTP_PORT"},
			Value:       defaultHTTPPort,
			Destination: &config.Port,
		},
		&cli.DurationFlag{
			Name:        "http-read-timeout",
			EnvVars:     []string{"HTTP_READ_TIMEOUT"},
			Value:       defaultHTTPRWTimeout,
			Destination: &config.ReadTimeout,
		},
		&cli.DurationFlag{
			Name:        "http-write-timeout",
			EnvVars:     []string{"HTTP_READ_WRITE_TIMEOUT"},
			Value:       defaultHTTPRWTimeout,
			Destination: &config.WriteTimeout,
		},
		&cli.IntFlag{
			Name:        "http-max-headers-megabytes",
			EnvVars:     []string{"HTTP_MAX_HEADER_MB"},
			Value:       defaultHTTPMaxHeaderMegabytes,
			Destination: &config.MaxHeaderMegabytes,
		},
		&cli.DurationFlag{
			Name:        "http-shutdown-delay",
			EnvVars:     []string{"HTTP_SHUTDOWN_DELAY"},
			Value:       defaultShutdownDuration,
			Destination: &config.ShutdownDelay,
		},
		&cli.DurationFlag{
			Name:        "http-shutdown-timeout",
			EnvVars:     []string{"HTTP_SHUTDOWN_TIMEOUT"},
			Value:       defaultShutdownTimeout,
			Destination: &config.ShutdownTimeout,
		},
	}

	return flags
}
