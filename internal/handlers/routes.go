package handlers

import (
	"github.com/labstack/echo/v4"
)

type ServerInterface interface {
	GetMainPage(echo.Context) error
}

func GetRoutes(si ServerInterface) *echo.Echo {
	e := echo.New()
	e.GET("/", si.GetMainPage)
	return e

}
