package google

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/novel/auth/config"
	"github.com/novel/auth/handler"
	"github.com/novel/auth/middleware"
	"github.com/novel/auth/usecase"
	"github.com/novel/auth/usecase/google"
	"golang.org/x/oauth2"
)

type googleAuthHandler struct {
	googleAuthUsecase usecase.IAuthUsecase
}

var googleConfig = oauth2.Config{
	ClientID:     config.Getenv("GOOGLE_ID"),
	ClientSecret: config.Getenv("GOOGLE_SECRET"),
	RedirectURL:  "http://localhost:12121/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://oauth2.googleapis.com/token",
	},
}

var instance *googleAuthHandler = nil

func New() *googleAuthHandler {
	if instance == nil {
		instance = &googleAuthHandler{
			googleAuthUsecase: google.New(),
		}
	}
	return instance
}

func (g *googleAuthHandler) Signin(ctx *middleware.Ctx) {
	state := handler.GenerateState(ctx.W)          // csrf token
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

	jwtToken, err := g.googleAuthUsecase.CreateUserToken(user)
	if err != nil {
		log.Println(err)
		http.Redirect(ctx.W, ctx.R, "/", http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(jwtToken)
	if err != nil {
		log.Println(err)
		http.Redirect(ctx.W, ctx.R, "/", http.StatusBadRequest)
		return
	}

	ctx.W.Write(response)
}
