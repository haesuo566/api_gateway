package naver

import (
	"context"
	"net/http"

	"github.com/novel/auth/config"
	"github.com/novel/auth/handler"
	"github.com/novel/auth/middleware"
	"github.com/novel/auth/usecase"
	"github.com/novel/auth/usecase/naver"
	"golang.org/x/oauth2"
)

type naverAuthHandler struct {
	naverAuthUsecase usecase.IAuthUsecase
}

var naverConfig = oauth2.Config{
	ClientID:     config.Getenv("NAVER_ID"),
	ClientSecret: config.Getenv("NAVER_SECRET"),
	RedirectURL:  "http://localhost:12121/auth/naver/callback",
	Scopes:       []string{"https://openapi.naver.com/v1/nid/me"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://nid.naver.com/oauth2.0/authorize",
		TokenURL: "https://nid.naver.com/oauth2.0/token",
	},
}

var instance *naverAuthHandler = nil

func New() *naverAuthHandler {
	if instance == nil {
		instance = &naverAuthHandler{
			naverAuthUsecase: naver.New(),
		}
	}
	return instance
}

func (n *naverAuthHandler) Signin(ctx *middleware.Ctx) {
	state := handler.GenerateState(ctx.W)
	url := naverConfig.AuthCodeURL(state)
	http.Redirect(ctx.W, ctx.R, url, http.StatusTemporaryRedirect)
}

func (n *naverAuthHandler) Callback(ctx *middleware.Ctx) {
	state, err := ctx.R.Cookie("state")
	if err != nil {
		http.Redirect(ctx.W, ctx.R, "/", http.StatusTemporaryRedirect)
		return
	}

	if state.Value != ctx.R.FormValue("state") {
		http.Redirect(ctx.W, ctx.R, "/", http.StatusTemporaryRedirect)
		return
	}

	code := ctx.R.FormValue("code")
	token, err := naverConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Redirect(ctx.W, ctx.R, "/", http.StatusTemporaryRedirect)
		return
	}

	n.naverAuthUsecase.GetUserInfo(token)
}
