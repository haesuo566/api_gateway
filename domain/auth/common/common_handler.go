package common

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/novel/api-gateway/util/jwt"
)

type commonHandler struct {
	usecase iCommonUsecase
}

var handlerInstance *commonHandler = nil

func newHandler(usecase iCommonUsecase) *commonHandler {
	if handlerInstance == nil {
		handlerInstance = &commonHandler{
			usecase: usecase,
		}
	}

	return handlerInstance
}

func (n *commonHandler) Login(ctx *fiber.Ctx) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	// login process
	user, err := n.usecase.login(email, password)
	if err != nil {
		return err
	}

	// create access/refresh token
	responseToken, err := jwt.GenerateResponseToken(user)
	if err != nil {
		return err
	}

	return ctx.JSON(responseToken)
}

func (n *commonHandler) Logout(ctx *fiber.Ctx) error {
	return nil
}

func (n *commonHandler) Signup(ctx *fiber.Ctx) error {
	username := ctx.FormValue("username")
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	if len(username) > 16 || len(username) <= 0 {
		return errors.New("username length is incorrect")
	}

	if len(email) > 64 || len(email) <= 0 {
		return errors.New("email length is incorrect")
	}

	if len(password) > 64 || len(password) <= 0 {
		return errors.New("password length is incorrect")
	}

	// 비밀번호 체크는 나중에 생각나면 다시하자고
	emailRegex, err := regexp.Compile(`[a-zA-Z0-9]+@[a-zA-Z0-9]+((\.[a-zA-Z0-9]+){1,5})`)
	if err != nil {
		return err
	}

	if !emailRegex.MatchString(email) {
		return errors.New("email not formmated")
	}

	if err := n.usecase.singup(username, email, password); err != nil {
		return err
	}

	url := fmt.Sprintf("/auth/common/login?email=%s&password=%s", email, password)
	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}

func (n *commonHandler) RefreshToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.FormValue("refreshToken")

	if err := jwt.ValidateToken(refreshToken); err != nil {
		return err
	}

	accessToken, err := jwt.GenerateAccessToken()
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"access_token": accessToken,
	})
}
