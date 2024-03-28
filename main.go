package main

import (
	"log"

	"github.com/novel/auth/handler"
	"github.com/novel/auth/middleware"
)

// '/' 경로는 Api Gateway 추가예정
func main() {
	app := middleware.New()

	// Handlers
	authHandler := handler.NewAuthHandler()
	googleAuthHandler := handler.NewGoogleAuthHandler()

	app.Get("/signin", authHandler.Signin)

	app.Get("/auth/google", googleAuthHandler.Signin)
	app.Get("/auth/google/callback", googleAuthHandler.Callback)

	if err := app.Run(":12121"); err != nil {
		log.Println(err)
	}
}
