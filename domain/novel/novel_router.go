package novel

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

	group := app.Group("/auth/novel")
	group.Post("/login", handler.Login)
	group.Post("/signup", handler.Signup)
}
