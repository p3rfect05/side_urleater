package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service interface {
}

type Handlers struct {
	Service Service
}

func (h *Handlers) GetMainPage(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"hello": "goodbye",
	})
}
