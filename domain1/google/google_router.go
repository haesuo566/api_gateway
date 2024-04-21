package google

import "github.com/labstack/echo/v4"

func SetRouter(e *echo.Echo) {
	// dependency injection
	repository := NewRepository()
	usecase := NewUsecase(repository)
	handler := NewHandler(usecase)

	group := e.Group("/google")
	group.GET("/login", handler.Login)
	group.GET("/callback", handler.Callback)
}
