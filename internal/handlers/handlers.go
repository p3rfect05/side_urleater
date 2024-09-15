package handlers

import (
	"context"
	"fmt"
	"github.com/antonlindstrom/pgstore"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"urleater/internal/repository/postgresDB"
)

type Service interface {
	LoginUser(ctx context.Context, email string, password string) error
	RegisterUser(ctx context.Context, email string, password string) error
	CreateShortLink(ctx context.Context, shortLink string, longLink string, userEmail string) (*postgresDB.Link, error)
	UpdateUserShortLinks(ctx context.Context, email string, deltaLinks int) (*postgresDB.User, error)
	GetAllUserShortLinks(ctx context.Context, email string) ([]postgresDB.Link, *postgresDB.User, error)
	GetSubscriptions(ctx context.Context) ([]postgresDB.Subscription, error)
	GetShortLink(ctx context.Context, shortLink string) (*postgresDB.Link, error)
	GetUser(ctx context.Context, email string) (*postgresDB.User, error)
}

// TODO populate
var domain string = "http://localhost:8080"

type Handlers struct {
	Service Service
	Store   *pgstore.PGStore
}

func retrieveEmailFromSession(c echo.Context, store *pgstore.PGStore) (string, error) {
	session, err := store.Get(c.Request(), "session_key")

	if err != nil {
		return "", fmt.Errorf("error getting session: %w", err)
	}

	if _, ok := session.Values["email"]; !ok {
		return "", nil
	}
	res := session.Values["email"].(string)
	return res, nil
}

func (h *Handlers) GetMainPage(c echo.Context) error {
	email, err := retrieveEmailFromSession(c, h.Store)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if email == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	return c.Render(http.StatusOK, "main_page.html", nil)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *Handlers) PostLogin(c echo.Context) error {
	email, err := retrieveEmailFromSession(c, h.Store)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if email != "" {
		return c.JSON(http.StatusOK, echo.Map{
			"redirect_to": "/",
		})
	}

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

	err = h.Service.LoginUser(ctx, requestData.Email, requestData.Password)

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

	return c.JSON(http.StatusOK, echo.Map{
		"redirect_to": "/",
	})
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

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *Handlers) PostRegister(c echo.Context) error {
	email, err := retrieveEmailFromSession(c, h.Store)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if email != "" {
		return c.JSON(http.StatusOK, echo.Map{
			"redirect_to": "/",
		})
	}

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

	err = h.Service.RegisterUser(ctx, requestData.Email, requestData.Password)

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

	return c.JSON(http.StatusOK, echo.Map{
		"redirect_to": "/",
	})
}

type CreateShortLinkRequest struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url" validate:"required"`
}

type CreateShortLinkResponse struct {
	Link postgresDB.Link `json:"link"`
}

func (h *Handlers) CreateShortLink(c echo.Context) error {
	email, err := retrieveEmailFromSession(c, h.Store)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if email == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"redirect_to": "/login",
		})
	}

	ctx := c.Request().Context()

	requestData := new(CreateShortLinkRequest)

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if c.Echo().Validator != nil {
		if err := c.Validate(requestData); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	link, err := h.Service.CreateShortLink(ctx, requestData.ShortURL, requestData.LongURL, email)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, CreateShortLinkResponse{
		Link: *link,
	})
}

type GetUserShortLinksResponse struct {
	Links []postgresDB.Link `json:"links"`
	User  postgresDB.User   `json:"user"`
}

func (h *Handlers) GetUserShortLinks(c echo.Context) error {
	email, err := retrieveEmailFromSession(c, h.Store)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if email == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"redirect_to": "/login",
		})
	}

	ctx := c.Request().Context()

	links, user, err := h.Service.GetAllUserShortLinks(ctx, email)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, GetUserShortLinksResponse{
		Links: links,
		User:  *user,
	})
}

func (h *Handlers) GetLoginPage(c echo.Context) error {
	email, err := retrieveEmailFromSession(c, h.Store)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if email != "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	return c.Render(http.StatusOK, "login_page.html", nil)
}

func (h *Handlers) GetRegisterPage(c echo.Context) error {
	email, err := retrieveEmailFromSession(c, h.Store)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if email != "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	return c.Render(http.StatusOK, "register_page.html", nil)
}

type UpdateUserShortLinksRequest struct {
	Email      string `json:"email" validate:"required"`
	DeltaLinks int    `json:"delta_links" validate:"required"`
}

type UpdateUserShortLinksResponse struct {
	User postgresDB.User `json:"user"`
}

func (h *Handlers) UpdateUserShortLinks(c echo.Context) error {
	email, err := retrieveEmailFromSession(c, h.Store)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if email != "admin@admin.com" {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("user %s is not authorized to change links number", email))
	}

	ctx := c.Request().Context()

	requestData := new(UpdateUserShortLinksRequest)

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if c.Echo().Validator != nil {
		if err := c.Validate(requestData); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	user, err := h.Service.UpdateUserShortLinks(ctx, requestData.Email, requestData.DeltaLinks)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, UpdateUserShortLinksResponse{
		User: *user,
	})

}

func (h *Handlers) GetCreateShortLink(c echo.Context) error {
	email, err := retrieveEmailFromSession(c, h.Store)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if email == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	return c.Render(http.StatusOK, "create_link_page.html", nil)
}

func (h *Handlers) GetShortLink(c echo.Context) error {
	ctx := c.Request().Context()

	shortLink := c.Param("short_link")

	link, err := h.Service.GetShortLink(ctx, shortLink)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusMovedPermanently, link.LongUrl)
}

type GetSubscriptionsResponse struct {
	Subscriptions []postgresDB.Subscription `json:"subscriptions"`
}

func (h *Handlers) GetSubscriptions(c echo.Context) error {
	email, err := retrieveEmailFromSession(c, h.Store)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if email == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"redirect_to": "/login",
		})
	}

	ctx := c.Request().Context()

	subscriptions, err := h.Service.GetSubscriptions(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, GetSubscriptionsResponse{
		Subscriptions: subscriptions,
	})
}

func (h *Handlers) GetSubscriptionsPage(c echo.Context) error {
	email, err := retrieveEmailFromSession(c, h.Store)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if email == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	return c.Render(http.StatusOK, "subscriptions.html", nil)

}

type GetUserResponse struct {
	User postgresDB.User `json:"user"`
}

func (h *Handlers) GetUser(c echo.Context) error {
	email, err := retrieveEmailFromSession(c, h.Store)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if email == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"redirect_to": "/login",
		})
	}

	ctx := c.Request().Context()

	user, err := h.Service.GetUser(ctx, email)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, GetUserResponse{
		User: *user,
	})

}
