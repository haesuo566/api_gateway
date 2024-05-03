package google

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

type GoogleHandler struct {
	usecase IGoogleUsecase
}

var googleConfig oauth2.Config

func init() {
	config.LoadEnv()
	googleConfig = oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_ID"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		RedirectURL:  "http://localhost:12121/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
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

func (g *GoogleHandler) Login(c echo.Context) error {
	state := GenerateToken(c)
	url := googleConfig.AuthCodeURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (g *GoogleHandler) Callback(c echo.Context) error {
	state, err := c.Cookie("state")
	if err != nil {
		c.Redirect(http.StatusBadRequest, "/google/login")
		return err
	}

	if c.FormValue("state") != state.Value {
		c.Redirect(http.StatusBadRequest, "/google/login")
		return errors.New("")
	}

	code := c.FormValue("code")
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		c.Redirect(http.StatusBadRequest, "/google/login")
		return err
	}

	user, err := g.usecase.GetUserInfo(token)
	if err != nil {
		c.Redirect(http.StatusBadRequest, "/google/login")
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
