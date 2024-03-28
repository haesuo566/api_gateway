package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/novel/auth/middleware"
	"github.com/novel/auth/usecase"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleAuthHandler struct {
	googleAuthUsecase usecase.IGoogleAuthUsecase
}

// .env file 로 뺌
var googleConfig = oauth2.Config{
	ClientID:     "1014459614066-945esdhcqevf8u9une9i0b7bvofsihld.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:12121/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func NewGoogleAuthHandler() *googleAuthHandler {
	return &googleAuthHandler{
		googleAuthUsecase: usecase.NewGoogleAuthUsecase(),
	}
}

func (g *googleAuthHandler) Signin(ctx *middleware.Ctx) {
	state := generateState(ctx.W)                  // csrf token
	redirectUrl := googleConfig.AuthCodeURL(state) // login redirect url -> ex) https://google.auth?state="asdasd"&code="asd"
	http.Redirect(ctx.W, ctx.R, redirectUrl, http.StatusTemporaryRedirect)
}

func (g *googleAuthHandler) Callback(ctx *middleware.Ctx) {
	state, err := ctx.R.Cookie("state")
	if err != nil {
		log.Println(err)
		http.Redirect(ctx.W, ctx.R, "/", http.StatusBadRequest)
		return
	}

	if ctx.R.FormValue("state") != state.Value {
		http.Redirect(ctx.W, ctx.R, "/", http.StatusBadRequest)
		return
	}

	code := ctx.R.FormValue("code")
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println(err)
		http.Redirect(ctx.W, ctx.R, "/", http.StatusBadRequest)
		return
	}

	g.googleAuthUsecase.GetUserInfo(token)
}
