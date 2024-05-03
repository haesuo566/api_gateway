package router

import (
	"github.com/labstack/echo/v4"
	"github.com/novel/auth/domain/google"
	"github.com/novel/auth/domain/naver"
)

func SetRouter(e *echo.Echo) {
	google.SetRouter(e)
	naver.SetRouter(e)
}
