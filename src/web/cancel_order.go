package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/buurzx/cryptoexchange/src/entities"
	"github.com/labstack/echo/v4"
)

type cancelOrder struct {
	orderbooksRepo OrderbooksRepoIface
	orderRepo      OrderRepoIface
	logger         Logger
}

func NewCancelOrderHandler(orderbookRepo OrderbooksRepoIface, orderRepo OrderRepoIface, log Logger) *cancelOrder {
	return &cancelOrder{
		orderRepo:      orderRepo,
		orderbooksRepo: orderbookRepo,
		logger:         log,
	}
}

func (co *cancelOrder) Handle(c echo.Context) error {
	orderbookIDStr := c.Param("orderbook_id")
	idStr := c.Param("id")

	orderbookID, err := strconv.Atoi(orderbookIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"errors": []string{err.Error()}})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"errors": []string{err.Error()}})
	}

	order, err := co.orderRepo.FindByID(int64(id))
	if err != nil && errors.Is(entities.ErrNotFound, err) {
		return c.JSON(http.StatusNotFound, map[string]any{"errors": notFoundErrorResponse})
	}

	if err != nil {
		co.logger.Errorf("failed to find order %w", err)
		return c.JSON(http.StatusInternalServerError, map[string]any{"errors": []string{"Something went wrong"}})
	}

	ob, err := co.orderbooksRepo.FindByID(int64(orderbookID))
	if err != nil && errors.Is(entities.ErrNotFound, err) {
		co.logger.Errorf("failed to find orderbook %w", err)
		return c.JSON(http.StatusNotFound, map[string]any{"errors": notFoundErrorResponse})
	}

	ob.CancelOrder(order)

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}
