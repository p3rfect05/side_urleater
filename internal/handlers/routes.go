package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
)

type ServerInterface interface {
	GetMainPage(echo.Context) error
	PostLogin(c echo.Context) error
	PostRegister(c echo.Context) error
	GetLogout(c echo.Context) error
	CreateShortLink(c echo.Context) error
	UpdateUserShortLinks(c echo.Context) error
	GetRegisterPage(c echo.Context) error
	GetLoginPage(c echo.Context) error
	GetUserShortLinks(c echo.Context) error
	GetCreateShortLink(c echo.Context) error
	GetShortLink(c echo.Context) error
	GetSubscriptions(c echo.Context) error
	GetSubscriptionsPage(c echo.Context) error
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func GetRoutes(si ServerInterface) *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORS())

	e.Use(middleware.Static("/static"))

	t := &Template{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}

	e.Renderer = t

	e.GET("/", si.GetMainPage)
	e.GET("/login", si.GetLoginPage)
	e.GET("/register", si.GetRegisterPage)
	e.POST("/login", si.PostLogin)
	e.POST("/register", si.PostRegister)
	e.GET("/logout", si.GetLogout)
	e.POST("/create_link", si.CreateShortLink)
	e.GET("/create_link", si.GetCreateShortLink)
	e.GET("/:short_link", si.GetShortLink)
	e.GET("/subscriptions", si.GetSubscriptionsPage)
	e.GET("/get_subscriptions", si.GetSubscriptions)

	return e

}
