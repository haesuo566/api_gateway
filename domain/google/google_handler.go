package google

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

type GoogleHandler struct {
	usecase IGoogleUsecase
}

var googleConfig oauth2.Config

func init() {
	config.LoadEnv()
	googleConfig = oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_ID"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		RedirectURL:  "http://localhost:12121/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
}

var handlerInstance *GoogleHandler = nil

func NewHandler(usecase IGoogleUsecase) *GoogleHandler {
	if handlerInstance == nil {
		handlerInstance = &GoogleHandler{
			usecase: usecase,
		}
	}
	return handlerInstance
}

func (g *GoogleHandler) Login(ctx *fiber.Ctx) error {
	state := GenerateToken(ctx)
	url := googleConfig.AuthCodeURL(state)
	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}

func (g *GoogleHandler) Callback(ctx *fiber.Ctx) error {
	state := ctx.Cookies("state")
	if ctx.FormValue("state") != state {
		ctx.Redirect("/google/login", http.StatusBadRequest)
		return errors.New("")
	}

	code := ctx.FormValue("code")
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.Redirect("/google/login", http.StatusBadRequest)
		return err
	}

	user, err := g.usecase.GetUserInfo(token)
	if err != nil {
		ctx.Redirect("/google/login", http.StatusBadRequest)
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
