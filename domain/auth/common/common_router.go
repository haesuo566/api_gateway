package common

import (
	"github.com/gofiber/fiber/v2"
	"github.com/novel/api-gateway/db"
	"github.com/novel/api-gateway/entity/user"
)

func SetRouter(app *fiber.App) {
	// dependency injection
	db := db.New()
	userRepository := user.NewRepository(db)
	usecase := newUsecase(userRepository)
	handler := newHandler(usecase)

	group := app.Group("/auth/common")
	group.Post("/login", handler.Login)
	group.Post("/signup", handler.Signup)
	group.Post("/refresh-token", handler.RefreshToken)
}
