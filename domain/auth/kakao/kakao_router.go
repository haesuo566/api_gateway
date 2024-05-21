package kakao

import (
	"github.com/labstack/echo/v4"
	"github.com/novel/api-gateway/db"
	"github.com/novel/api-gateway/entity/user"
)

func SetRouter(e *echo.Echo) {
	// dependency injection
	db := db.New()
	userRepository := user.NewRepository(db)
	usecase := NewUsecase(userRepository)
	handler := NewHandler(usecase)

	group := e.Group("/auth/kakao")
	group.GET("/login", handler.Login)
	group.GET("/callback", handler.Callback)
}
