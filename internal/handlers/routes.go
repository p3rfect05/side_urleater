package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerInterface interface {
	GetMainPage(echo.Context) error
	PostLogin(c echo.Context) error
	PostRegister(c echo.Context) error
	GetLogout(c echo.Context) error
	CreateShortLink(c echo.Context) error
}

func GetRoutes(si ServerInterface) *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())
	
	e.GET("/", si.GetMainPage)
	e.POST("/login", si.PostLogin)
	e.POST("/register", si.PostRegister)
	e.GET("/logout", si.GetLogout)
	e.POST("/create_link", si.CreateShortLink)

	return e

}
