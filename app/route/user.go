package route

import (
	"github.com/labstack/echo/v5"
)

type userHandler interface{
	ByName(*echo.Context) error
}

func userRoute(e *echo.Echo, uh userHandler) {
	u := e.Group("/users")
	u.GET("/:name", uh.ByName)
}
