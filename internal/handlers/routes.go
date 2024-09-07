package handlers

import (
	"github.com/labstack/echo/v4"
	"urleater/internal/validator"
)

type ServerInterface interface {
	GetMainPage(echo.Context) error
	PostLogin(c echo.Context) error
	PostRegister(c echo.Context) error
	GetLogout(c echo.Context) error
}

func GetRoutes(si ServerInterface) *echo.Echo {
	e := echo.New()
	httpValidator, err := validator.NewValidator()
	if err != nil {
		panic(err)
	}
	e.Validator = httpValidator
	e.GET("/", si.GetMainPage)
	e.POST("/login", si.PostLogin)
	e.POST("/register", si.PostRegister)
	e.GET("/logout", si.GetLogout)

	return e

}
