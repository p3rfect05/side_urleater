package handlers

import (
	"context"
	"fmt"
	"github.com/antonlindstrom/pgstore"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type Service interface {
	LoginUser(ctx context.Context, email string, password string) error
	RegisterUser(ctx context.Context, email string, password string) error
}

type Handlers struct {
	Service Service
	Store   *pgstore.PGStore
}

func (h *Handlers) GetMainPage(c echo.Context) error {
	session, err := h.Store.Get(c.Request(), "session_key")
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handlers) PostLogin(c echo.Context) error {
	ctx := c.Request().Context()

	requestData := new(LoginRequest)

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if c.Echo().Validator != nil {
		if err := c.Validate(requestData); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	err := h.Service.LoginUser(ctx, requestData.Email, requestData.Password)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	session, err := h.Store.Get(c.Request(), "session_key")
	if err != nil {
		log.Printf("Error getting session: %v\n", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	session.Values["email"] = requestData.Email

	if err = session.Save(c.Request(), c.Response()); err != nil {
		log.Printf("Error saving session: %v\n", err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}
	return c.JSON(http.StatusOK, echo.Map{})
}

func (h *Handlers) GetLogout(c echo.Context) error {
	session, err := h.Store.Get(c.Request(), "session_key")

	if err != nil {
		log.Printf("Error getting session: %v\n", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	session.Options.MaxAge = -1
	if err = session.Save(c.Request(), c.Response()); err != nil {
		log.Printf("Error saving session: %v\n", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{})
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handlers) PostRegister(c echo.Context) error {
	ctx := c.Request().Context()

	requestData := new(RegisterRequest)

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if c.Echo().Validator != nil {
		if err := c.Validate(requestData); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	err := h.Service.RegisterUser(ctx, requestData.Email, requestData.Password)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	session, err := h.Store.Get(c.Request(), "session_key")
	if err != nil {
		log.Printf("Error getting session: %v\n", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	session.Values["email"] = requestData.Email

	if err = session.Save(c.Request(), c.Response()); err != nil {
		log.Printf("Error saving session: %v\n", err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}
	return c.JSON(http.StatusOK, echo.Map{})
}
