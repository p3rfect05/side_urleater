package handlers

import (
	"fmt"
	"github.com/antonlindstrom/pgstore"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service interface {
}

type Handlers struct {
	Service Service
	Store   *pgstore.PGStore
}

func (h *Handlers) GetMainPage(c echo.Context) error {
	session, err := h.Store.Get(c.Request(), "session_id")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Printf("%v\n", session.Values)
	if len(session.Values) == 0 {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"hello": "goodbye",
	})
}
