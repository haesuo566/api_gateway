package google

import (
	"github.com/gofiber/fiber/v2"
	"github.com/novel/api-gateway/db"
	"github.com/novel/api-gateway/entity/user"
)

func SetRouter(app *fiber.App) {
	// dependency injection
	db := db.New()
	userRepository := user.NewRepository(db)
	usecase := NewUsecase(userRepository)
	handler := NewHandler(usecase)

	group := app.Group("/auth/google")
	group.Get("/login", handler.Login)
	group.Get("/callback", handler.Callback)
}
