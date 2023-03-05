package web

import (
	"context"
	"fmt"
	"time"

	"github.com/buurzx/cryptoexchange/src/app"
	configs "github.com/buurzx/cryptoexchange/src/app/configs"
	"github.com/buurzx/cryptoexchange/src/entities"
	"github.com/buurzx/cryptoexchange/src/repositories"
	"github.com/labstack/echo/v4"
)

type API struct {
	echo   *echo.Echo
	config *configs.WebConfig
	logger Logger
}

type Logger interface {
	Info(i ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type OrderbooksRepoIface interface {
	FindByMarket(string) (*entities.Orderbook, error)
	FindByID(int64) (*entities.Orderbook, error)
}

type OrderRepoIface interface {
	FindByID(int64) (*entities.Order, error)
}

var notFoundErrorResponse = []string{"record not found"}

const pluginName = "web"

func New() *API {
	return &API{
		echo: echo.New(),
	}
}

func (s *API) Setup(a *app.Application) error {
	s.config = a.Config.HTTP

	s.echo.Server.ReadTimeout = s.config.ReadTimeout
	s.echo.Server.WriteTimeout = s.config.WriteTimeout
	s.echo.Server.MaxHeaderBytes = s.config.MaxHeaderMegabytes
	s.echo.Server.Addr = s.config.Port

	s.logger = a.Logger

	// repos
	// TODO: remove repo dependency
	repos := repositories.NewRepos()
	err := repos.Setup(a)
	if err != nil {
		return fmt.Errorf("failed setup repo %w", err)
	}

	// handlers
	healthcheckHandler := NewHealthcheck()
	placeOrderHandler := NewPlaceOrderHandler(repos.Orderbook)
	getOrderHandler := NewGetOrderHandler(repos.Orderbook)
	cancelOrderHandler := NewCancelOrderHandler(repos.Orderbook, repos.Order, s.logger)

	// routes
	s.echo.GET("/health", healthcheckHandler.Handle)
	s.echo.POST("/orders", placeOrderHandler.Handle)
	s.echo.GET("/orderbook/:market", getOrderHandler.Handle)
	s.echo.DELETE("/orderbooks/:orderbook_id/orders/:id", cancelOrderHandler.Handle)

	return nil
}

func (s *API) Run() error {
	s.logger.Info("Starting web server ...")

	return s.echo.Start(s.config.Port)
}

func (s *API) Stop() error {
	s.logger.Infof("wait %s, http server shutdown ...", s.config.ShutdownDelay)
	time.Sleep(s.config.ShutdownDelay)

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		s.config.ShutdownTimeout,
	)
	defer cancel()

	if err := s.echo.Server.Shutdown(shutdownCtx); err != nil {
		s.logger.Errorf("failed to close http server %s", err.Error())

		return err
	}

	s.logger.Info("closed http server successfully")

	return nil
}

func (s *API) Name() string {
	return pluginName
}
