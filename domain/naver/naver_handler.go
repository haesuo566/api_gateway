package naver

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
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
		RedirectURL:  "http://localhost:12121/naver/callback",
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

func (g *naverHandler) login(c echo.Context) error {
	state := GenerateToken(c)
	url := naverConfig.AuthCodeURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (g *naverHandler) callback(c echo.Context) error {
	state, err := c.Cookie("state")
	if err != nil {
		c.Redirect(http.StatusBadRequest, "/naver/login")
		return err
	}

	if c.FormValue("state") != state.Value {
		c.Redirect(http.StatusBadRequest, "/naver/login")
		return errors.New("")
	}

	code := c.FormValue("code")
	token, err := naverConfig.Exchange(context.Background(), code)
	if err != nil {
		c.Redirect(http.StatusBadRequest, "/naver/login")
		return err
	}

	user, err := g.usecase.getUserInfo(token)
	if err != nil {
		c.Redirect(http.StatusBadRequest, "/naver/login")
		return err
	}

	responseToken, err := jwt.GenerateResponseToken(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, responseToken)
}

func GenerateToken(c echo.Context) string {
	expires := time.Now().Add(time.Hour * 24)

	data := make([]byte, 16)
	rand.Read(data)
	csrfToken := base64.URLEncoding.EncodeToString(data)

	cookie := &http.Cookie{}
	cookie.Name = "state"
	cookie.Value = csrfToken
	cookie.Expires = expires

	c.SetCookie(cookie)
	return csrfToken
}
