package naver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/novel/auth/db"
	"github.com/novel/auth/entity/user"
)

func SetRouter(app *fiber.App) {
	// dependency injection
	db := db.New()
	userRepository := user.NewRepository(db)
	usecase := newUsecase(userRepository)
	handler := newHandler(usecase)

	group := app.Group("/auth/naver")
	group.Get("/login", handler.Login)
	group.Get("/callback", handler.Callback)
}
