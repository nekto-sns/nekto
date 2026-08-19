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
		if internalErr := errors.Unwrap(err); internalErr != nil && appErr.Code != http.StatusNotFound && appErr.Code != http.StatusBadRequest {

			slog.Error("Request processing failed", "err", internalErr)
		}
		_ = c.JSON(appErr.Code, errorResponse{
			Message: appErr.Message,
		})
		return
	}

	var sc echo.HTTPStatusCoder
	if errors.As(err, &sc) {
		if code := sc.StatusCode(); code != 0 {
			if code != http.StatusNotFound {
				slog.Warn("HTTP error", "err", err.Error())
			}
			_ = c.JSON(code, errorResponse{
				Message: err.Error(),
			})
			return
		}
	}


	_ = c.JSON(http.StatusInternalServerError, errorResponse{
		Message: "Internal server error",
	})
	return
}
