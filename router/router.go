package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/novel/api-gateway/domain/auth/common"
	"github.com/novel/api-gateway/domain/auth/google"
	"github.com/novel/api-gateway/domain/auth/naver"
)

func SetRouter(app *fiber.App) {
	google.SetRouter(app)
	naver.SetRouter(app)
	common.SetRouter(app)
	// kakao.SetRouter(app) // kakao 는 심사를 받아야 email 정보를 얻을 수 있어서 배제함
}
