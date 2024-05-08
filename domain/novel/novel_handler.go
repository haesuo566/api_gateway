package novel

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/novel/auth/util/jwt"
)

type novelHandler struct {
	usecase iNovelUsecase
}

var handlerInstance *novelHandler = nil

func newHandler(usecase iNovelUsecase) *novelHandler {
	if handlerInstance == nil {
		handlerInstance = &novelHandler{
			usecase: usecase,
		}
	}

	return handlerInstance
}

func (n *novelHandler) login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

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

	return c.JSON(http.StatusOK, responseToken)
}

func (n *novelHandler) logout() {

}

func (n *novelHandler) signup(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

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
	emailRegex, err := regexp.Compile(`[a-zA-Z0-9]+@[a-zA-Z0-9]+((\\.[a-zA-Z0-9]+){1,5})`)
	if err != nil {
		return err
	}

	if !emailRegex.MatchString(email) {
		return errors.New("email not formmated")
	}

	if err := n.usecase.singup(username, email, password); err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/auth/novel/login")
}
