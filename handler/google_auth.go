package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/novel/auth/ancho"
	"golang.org/x/oauth2"
)

type googleAuthHandler struct {
}

// .env file 로 뺌
var googleConfig = oauth2.Config{
	ClientID:     "1014459614066-945esdhcqevf8u9une9i0b7bvofsihld.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://oauth2.googleapis.com/token",
	},
}

// usecase 까지는 나눠야 하는데...
func NewGoogleAuthHandler() *googleAuthHandler {
	return &googleAuthHandler{}
}

func (g *googleAuthHandler) Signin(ctx *ancho.Ctx) {
	state := generateState(ctx.W)                  // csrf token
	redirectUrl := googleConfig.AuthCodeURL(state) // login redirect url -> ex) https://google.auth?state="asdasd"&code="asd"
	http.Redirect(ctx.W, ctx.R, redirectUrl, http.StatusTemporaryRedirect)
}

func (g *googleAuthHandler) Callback(ctx *ancho.Ctx) {
	state, err := ctx.R.Cookie("state")
	if err != nil {
		http.Redirect(ctx.W, ctx.R, "/", http.StatusBadRequest)
		return
	}

	if ctx.R.URL.Query().Get("state") != state.Value {
		http.Redirect(ctx.W, ctx.R, "/", http.StatusBadRequest)
		return
	}

	token, err := googleConfig.Exchange(context.Background(), ctx.R.URL.Query().Get("code"))
	if err != nil {
		http.Redirect(ctx.W, ctx.R, "/", http.StatusBadRequest)
		return
	}

	log.Println(token)
}
