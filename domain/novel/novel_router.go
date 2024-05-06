package novel

import (
	"github.com/labstack/echo/v4"
	"github.com/novel/auth/db"
	"github.com/novel/auth/entity/user"
)

func SetRouter(e *echo.Echo) {
	// dependency injection
	db := db.New()
	userRepository := user.NewRepository(db)
	usecase := newUsecase(userRepository)
	handler := newHandler(usecase)

	group := e.Group("/novel")
	group.GET("/login", handler.login)
	group.GET("/signup", handler.signup)
}
