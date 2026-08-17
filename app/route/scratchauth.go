package route

import (
	"github.com/labstack/echo/v5"
)

type scratchAuthHandler interface{
	LoginRedirect(*echo.Context) error
	LoginCallback(*echo.Context) error
}

func saRoute(e *echo.Echo, sh scratchAuthHandler) {
	e.Group("/auth")
	e.GET("/login", sh.LoginRedirect)
	e.GET("/login/callback", sh.LoginCallback)
}
