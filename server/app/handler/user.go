package handler

import (
	"time"
	"errors"
	"context"
	"net/http"
	"github.com/labstack/echo/v5"

	"github.com/nekto-sns/nekto/server/app/model"
)

type userService interface{
	ByName(context.Context, string) (*model.User, error)
}

type userHandler struct{
	svc userService
}

type userResponse struct{
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Bio         string    `json:"bio"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewUser(svc userService) (*userHandler) {
	return &userHandler{
		svc: svc,
	}
}

func (h *userHandler) ByName(c *echo.Context) error {
	name := c.Param("name")

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	user, err := h.svc.ByName(ctx, name)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user").Wrap(err)
	}
	return c.JSON(http.StatusOK, &userResponse{
		ID: user.ID,
		Name: user.Name,
		DisplayName: user.DisplayName,
		Bio: user.Bio,
		CreatedAt: user.CreatedAt,
	})
}
