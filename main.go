package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/novel/api-gateway/router"
)

/*
logout(나중에 구현) 은 redis 로 처리. 나중에 캐시데이터도 redis 에 넣는걸로
jwt, oauth 보안 블로그 -> https://puleugo.tistory.com/139
*/

// '/' 경로는 Api Gateway 추가예정
func main() {
	app := fiber.New()

	app.Use(logger.New())
	// app.Use(recover.New())

	router.SetRouter(app)

	app.Listen(":12121")
}
