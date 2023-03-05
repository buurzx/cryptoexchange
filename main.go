package main

import (
	"context"
	"os"

	"github.com/buurzx/cryptoexchange/src/cmd/app"
	"github.com/urfave/cli/v2"
)

func main() {
	newApp := &cli.App{
		Commands: []*cli.Command{
			app.BuildCmd(),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := newApp.RunContext(ctx, os.Args); err != nil {
		panic(err)
	}
}
