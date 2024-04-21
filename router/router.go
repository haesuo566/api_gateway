package router

import (
	"github.com/labstack/echo/v4"
	"github.com/novel/auth/domain1/google"
)

func SetRouter(e *echo.Echo) {
	google.SetRouter(e)
}
