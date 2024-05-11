package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/novel/auth/domain/google"
	"github.com/novel/auth/domain/naver"
)

func SetRouter(app *fiber.App) {
	google.SetRouter(app)
	naver.SetRouter(app)
	// novel.SetRouter(app)
	// kakao.SetRouter(app) // kakao 는 심사를 받아야 email 정보를 얻을 수 있어서 배제함
}
