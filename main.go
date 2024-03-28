package main

import (
	"log"

	"github.com/novel/auth/ancho"
	"github.com/novel/auth/handler"
)

func main() {
	app := ancho.New()

	authHandler := handler.NewAuthHandler()
	googleAuthHandler := handler.NewGoogleAuthHandler()

	app.Get("/signin", authHandler.Signin)

	app.Get("/auth/google", googleAuthHandler.Signin)
	app.Get("/auth/google/callback", googleAuthHandler.Callback)

	if err := app.Run(":8080"); err != nil {
		log.Println(err)
	}
}
