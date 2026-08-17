package route

import (
	"github.com/labstack/echo/v5"
)

func Setup(e *echo.Echo, uh userHandler, sh scratchAuthHandler) {
	userRoute(e, uh)
	saRoute(e, sh)
}
