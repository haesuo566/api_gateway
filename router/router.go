package router

import (
	"github.com/labstack/echo/v4"
	"github.com/novel/auth/domain/google"
	"github.com/novel/auth/domain/kakao"
	"github.com/novel/auth/domain/naver"
)

func SetRouter(e *echo.Echo) {
	google.SetRouter(e)
	naver.SetRouter(e)
	kakao.SetRouter(e) // kakao 는 심사를 받아야 email 정보를 얻을 수 있어서 배제함
}
