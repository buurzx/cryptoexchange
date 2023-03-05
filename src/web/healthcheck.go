package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type healthcheck struct {
}

func NewHealthcheck() healthcheck {
	return healthcheck{}
}

func (h healthcheck) Handle(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"health": "ok"})
}
