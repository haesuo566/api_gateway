package kakao

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

type KakaoHandler struct {
	usecase IKakaoUsecase
}

var naverConfig oauth2.Config

func init() {
	config.LoadEnv()
	naverConfig = oauth2.Config{
		ClientID:     os.Getenv("KAKAO_ID"),
		ClientSecret: os.Getenv("KAKAO_SECRET"),
		RedirectURL:  "http://localhost:12121/kakao/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://kauth.kakao.com/oauth/authorize",
			TokenURL: "https://kauth.kakao.com/oauth/token",
		},
	}
}

var handlerInstance *KakaoHandler = nil

func NewHandler(usecase IKakaoUsecase) *KakaoHandler {
	if handlerInstance == nil {
		handlerInstance = &KakaoHandler{
			usecase: usecase,
		}
	}
	return handlerInstance
}

func (g *KakaoHandler) Login(c echo.Context) error {
	state := GenerateToken(c)
	url := naverConfig.AuthCodeURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (g *KakaoHandler) Callback(c echo.Context) error {
	state, err := c.Cookie("state")
	if err != nil {
		c.Redirect(http.StatusBadRequest, "/kakao/login")
		return err
	}

	if c.FormValue("state") != state.Value {
		c.Redirect(http.StatusBadRequest, "/kakao/login")
		return errors.New("")
	}

	code := c.FormValue("code")
	token, err := naverConfig.Exchange(context.Background(), code)
	if err != nil {
		c.Redirect(http.StatusBadRequest, "/kakao/login")
		return err
	}

	user, err := g.usecase.GetUserInfo(token)
	if err != nil {
		c.Redirect(http.StatusBadRequest, "/kakao/login")
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
