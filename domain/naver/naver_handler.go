package naver

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/novel/auth/config"
	"github.com/novel/auth/util/jwt"
	"golang.org/x/oauth2"
)

type naverHandler struct {
	usecase iNaverUsecase
}

var naverConfig oauth2.Config

func init() {
	config.LoadEnv()
	naverConfig = oauth2.Config{
		ClientID:     os.Getenv("NAVER_ID"),
		ClientSecret: os.Getenv("NAVER_SECRET"),
		RedirectURL:  "http://localhost:12121/auth/naver/callback",
		Scopes:       []string{"https://openapi.naver.com/v1/nid/me"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://nid.naver.com/oauth2.0/authorize",
			TokenURL: "https://nid.naver.com/oauth2.0/token",
		},
	}
}

var handlerInstance *naverHandler = nil

func newHandler(usecase iNaverUsecase) *naverHandler {
	if handlerInstance == nil {
		handlerInstance = &naverHandler{
			usecase: usecase,
		}
	}
	return handlerInstance
}

func (g *naverHandler) Login(ctx *fiber.Ctx) error {
	state := GenerateToken(ctx)
	url := naverConfig.AuthCodeURL(state)
	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}

func (g *naverHandler) Callback(ctx *fiber.Ctx) error {
	state := ctx.Cookies("state")
	if ctx.FormValue("state") != state {
		ctx.Redirect("/naver/login", http.StatusBadRequest)
		return errors.New("state is not equals")
	}

	code := ctx.FormValue("code")
	token, err := naverConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.Redirect("/naver/login", http.StatusBadRequest)
		return err
	}

	user, err := g.usecase.getUserInfo(token)
	if err != nil {
		ctx.Redirect("/naver/login", http.StatusBadRequest)
		return err
	}

	responseToken, err := jwt.GenerateResponseToken(user)
	if err != nil {
		return err
	}

	return ctx.JSON(responseToken)
}

func GenerateToken(ctx *fiber.Ctx) string {
	data := make([]byte, 16)
	rand.Read(data)
	csrfToken := base64.URLEncoding.EncodeToString(data)

	ctx.Cookie(&fiber.Cookie{
		Name:    "state",
		Value:   csrfToken,
		Expires: time.Now().Add(time.Hour * 24),
	})
	return csrfToken
}
