package naver

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

type NaverHandler struct {
	naverUsecase impl.AuthImpl
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

var handlerInstance *NaverHandler = nil

func NewHandler() *NaverHandler {
	if handlerInstance == nil {
		handlerInstance = &NaverHandler{
			naverUsecase: NewUsecase(),
		}
	}
	return handlerInstance
}

func (n *NaverHandler) Signin(ctx *middleware.Ctx) {
	state := auth.GenerateState(ctx.W)
	url := naverConfig.AuthCodeURL(state)
	http.Redirect(ctx.W, ctx.R, url, http.StatusTemporaryRedirect)
}

func (n *NaverHandler) Callback(ctx *middleware.Ctx) {
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

	user, err := n.naverUsecase.GetUserInfo(token)
	if err != nil {
		log.Println(err)
		return
	}

	responseToken, err := n.naverUsecase.CreateUserToken(user)
	if err != nil {
		log.Println(err)
		return
	}

	data, err := json.Marshal(responseToken)
	if err != nil {
		log.Println(err)
		return
	}

	ctx.W.Header().Set("Content-Type", "application/json")
	ctx.W.Write(data)
}
