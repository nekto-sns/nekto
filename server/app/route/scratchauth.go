package route

import (
	"github.com/labstack/echo/v5"
)

type scratchAuthHandler interface{
	LoginRedirect(*echo.Context) error
	LoginCallback(*echo.Context) error
}

func saRoute(e *echo.Echo, sh scratchAuthHandler) {
	s := e.Group("/auth/scratch")
	s.GET("/login", sh.LoginRedirect)
	s.GET("/login/callback", sh.LoginCallback)
}
