package handler

import (
	"context"
	"encoding/base64"
	"net/url"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/nekto-sns/nekto-server/app/model"
)

type authService interface{
	LoginCallback(context.Context, string) (string, bool, error)
}

type authHandler struct{
	svc authService
	loginRedirectURL string
}

func NewAuthHandler(svc authService, scratchAuthURL, loginCallbackURL string) (*authHandler) {
	base64URL := base64.URLEncoding.EncodeToString([]byte(loginCallbackURL))
	redirect, _ := url.JoinPath(scratchAuthURL, "/auth")

	return &authHandler{
		svc: svc,
		loginRedirectURL: redirect + "?redirect=" + base64URL,
	}
}

func (h *authHandler) LoginRedirect(c *echo.Context) error {
	return c.Redirect(303, h.loginRedirectURL)
}

func (h *authHandler) LoginCallback(c *echo.Context) error {
	privateCode := c.QueryParam("privateCode")
	if privateCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "privateCode not found")
	}

	req := c.Request().Context()

	userID, isValid, err := h.svc.LoginCallback(req, privateCode)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found").Wrap(err)
		} else if errors.Is(err, model.ErrForbidden) {
			return echo.NewHTTPError(http.StatusForbidden, "Unauthorized").Wrap(err)
		}
	}

	if !isValid {
		return echo.NewHTTPError(http.StatusForbidden, "Unauthorized")
	}

	return c.String(http.StatusOK, userID + " authorized")
}
