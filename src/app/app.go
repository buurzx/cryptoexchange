package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	configs "github.com/buurzx/cryptoexchange/src/app/configs"
	"github.com/buurzx/cryptoexchange/src/logger"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
)

type ErrorGroup struct {
	Group *errgroup.Group
	Ctx   context.Context
}

// Application holds configs and plugins to start
type Application struct {
	Config     *configs.Config
	Logger     *logger.Logger
	ErrorGroup ErrorGroup
	Plugins    map[string]Plugin
	plugins    []Plugin
}

type Plugin interface {
	Name() string
	Setup(*Application) error
	Run() error
	Stop() error
}

func New(ctx context.Context, cfg *configs.Config, logger *logger.Logger) (*Application, error) {
	group, ctx := errgroup.WithContext(ctx)
	errGroup := ErrorGroup{
		Group: group,
		Ctx:   ctx,
	}

	return &Application{
		Config:     cfg,
		Logger:     logger,
		ErrorGroup: errGroup,
		Plugins:    map[string]Plugin{},
	}, nil
}

// Register adds plugins to application for handling
// TODO: plugins is a slice
// Plugins - is a map name => index in slice
// create getByName
func (a *Application) Register(p Plugin) *Application {
	a.Plugins[p.Name()] = p
	a.plugins = append(a.plugins, p)

	return a
}

func (a *Application) Run(ctx context.Context) error {
	if err := a.setupPlugins(); err != nil {
		return fmt.Errorf("failed setup plugins: %w", err)
	}

	// Start application.
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	err := a.startPlugins()
	if err != nil {
		a.stopPlugins()
		return fmt.Errorf("failed start plugins: %w", err)
	}

	select {
	case <-exit:
		a.Logger.Info("stopping application")
	case <-a.ErrorGroup.Ctx.Done():
		a.Logger.Error("stopping application with error")
	case <-ctx.Done():
		a.Logger.Error("stopping application with context canceled")
	}

	signal.Stop(exit)
	a.Close(ctx)

	return nil
}

// Close closes connections and stops all plugins
func (a *Application) Close(_ctx context.Context) {
	if err := a.stopPlugins(); err != nil {
		a.Logger.Errorf("failed stop plugins: %v", err)
	}
}

func (a *Application) setupPlugins() error {
	var err error

	for _, plugin := range a.plugins {
		a.Logger.Debugf("Setuping plugin %s ...", plugin.Name())

		if perr := plugin.Setup(a); perr != nil {
			a.Logger.Errorf("failed setup plugin %s, %w", plugin.Name(), perr)
			err = multierr.Append(err, perr)
		}
	}

	return err
}

func (a *Application) startPlugins() error {
	for _, plugin := range a.plugins {
		a.Logger.Debugf("Starting plugin %s ...", plugin.Name())
		a.startPlugin(plugin)
	}

	return nil
}

func (a *Application) startPlugin(p Plugin) {
	a.ErrorGroup.Group.Go(func() error {
		if err := p.Run(); err != nil {
			a.Logger.Errorf("failed run plugin %s, %w", p.Name(), err)
			return err
		}

		return nil
	})
}

func (a *Application) stopPlugins() error {
	var err error

	for _, plugin := range a.plugins {
		a.Logger.Debugf("Stopping plugin %s ...", plugin.Name())
		if perr := plugin.Stop(); perr != nil {
			a.Logger.Errorf("failed stop plugin %s, %w", plugin.Name(), perr)
			err = multierr.Append(err, perr)
		}
	}

	return err
}
