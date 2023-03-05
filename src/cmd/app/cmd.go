package app

import (
	"fmt"

	"github.com/buurzx/cryptoexchange/src/app"
	appconfig "github.com/buurzx/cryptoexchange/src/app/configs"
	"github.com/buurzx/cryptoexchange/src/logger"
	"github.com/buurzx/cryptoexchange/src/repositories"
	"github.com/buurzx/cryptoexchange/src/web"
	"github.com/urfave/cli/v2"
)

func BuildCmd() *cli.Command {
	cfg := appconfig.NewConfig()

	return &cli.Command{
		Name:        "http-server",
		Description: "http server",
		Flags:       cfg.BuildFlags(),
		Action: func(ctx *cli.Context) error {
			logger, err := logger.New(cfg.LogLevel)
			if err != nil {
				return fmt.Errorf("failed to initialize logger %w", err)
			}

			app, err := app.New(ctx.Context, cfg, logger)
			if err != nil {
				return fmt.Errorf("failed to initialize application %w", err)
			}

			app.Register(repositories.NewOrderbooksRepo()).
				Register(web.New())

			err = app.Run(ctx.Context)
			if err != nil {
				return fmt.Errorf("failed to run application %w", err)
			}

			return nil
		},
	}
}
