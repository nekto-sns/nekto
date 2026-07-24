package errorhandler

import (
	"log/slog"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

type errorResponse struct{
	Message string `json:"message"`
}

func ErrorHandler(c *echo.Context, err error) {
	var appErr *echo.HTTPError
	if errors.As(err, &appErr) {
		if internalErr := errors.Unwrap(err); internalErr != nil {
			slog.Error("Request processing failed", "err", internalErr)
		}
		_ = c.JSON(appErr.Code, errorResponse{
			Message: appErr.Message,
		})
		return
	}
	_ = c.JSON(http.StatusInternalServerError, errorResponse{
		Message: "Internal server error",
	})
	return
}
