package google

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/novel/auth/global/common/impl"
	"github.com/novel/auth/global/config"
	"github.com/novel/auth/global/middleware"
	"github.com/novel/auth/global/util/auth"
	"golang.org/x/oauth2"
)

type GoogleHandler struct {
	googleUsecase impl.AuthImpl
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

var handlerInstance *GoogleHandler = nil

func NewHandler() *GoogleHandler {
	if handlerInstance == nil {
		handlerInstance = &GoogleHandler{
			googleUsecase: NewUsecase(),
		}
	}
	return handlerInstance
}

func (g *GoogleHandler) Signin(ctx *middleware.Ctx) {
	state := auth.GenerateState(ctx.W)             // csrf token
	redirectUrl := googleConfig.AuthCodeURL(state) // login redirect url -> ex) https://google.auth?state="asdasd"&code="asd"
	http.Redirect(ctx.W, ctx.R, redirectUrl, http.StatusTemporaryRedirect)
}

func (g *GoogleHandler) Callback(ctx *middleware.Ctx) {
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

	user, err := g.googleUsecase.GetUserInfo(token)
	if err != nil {
		log.Println(err)
		http.Redirect(ctx.W, ctx.R, "/", http.StatusBadRequest)
		return
	}

	jwtToken, err := g.googleUsecase.CreateUserToken(user)
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

	ctx.W.Header().Set("Content-Type", "application/json")
	ctx.W.Write(response)
}
