package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/novel/auth/config"
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
	ClientID:     config.Getenv("GOOGLE_ID"),
	ClientSecret: config.Getenv("GOOGLE_SECRET"),
	RedirectURL:  "http://localhost:12121/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

var instance *googleAuthHandler = nil

func NewGoogleAuthHandler() *googleAuthHandler {
	if instance == nil {
		instance = &googleAuthHandler{
			googleAuthUsecase: usecase.NewGoogleAuthUsecase(),
		}
	}
	return instance
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

	user, err := g.googleAuthUsecase.GetUserInfo(token)
	if err != nil {
		log.Println(err)
		http.Redirect(ctx.W, ctx.R, "/", http.StatusBadRequest)
		return
	}

	// jwt token return
	log.Println(user)
}
