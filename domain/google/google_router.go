package google

import (
	"github.com/labstack/echo/v4"
	"github.com/novel/auth/db"
	"github.com/novel/auth/entity/user"
)

func SetRouter(e *echo.Echo) {
	// dependency injection
	db := db.New()
	userRepository := user.NewRepository(db)
	usecase := NewUsecase(userRepository)
	handler := NewHandler(usecase)

	group := e.Group("/auth/google")
	group.GET("/login", handler.Login)
	group.GET("/callback", handler.Callback)
}
