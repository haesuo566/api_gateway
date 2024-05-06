package novel

import (
	"errors"
	"regexp"

	"github.com/labstack/echo/v4"
)

type novelHandler struct {
	usecase INovelUsecase
}

var handlerInstance *novelHandler = nil

func newHandler(usecase INovelUsecase) *novelHandler {
	if handlerInstance == nil {
		handlerInstance = &novelHandler{
			usecase: usecase,
		}
	}

	return handlerInstance
}

func (n *novelHandler) login(c echo.Context) error {
	return nil
}

func (n *novelHandler) logout() {

}

func (n *novelHandler) signup(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return errors.New("")
	}

	// 이건 라이브러리 사용해야 할 것 같음
	regex, err := regexp.Compile("^(?=.*[a-zA-Z])(?=.*[0-9]).{8,25}$")
	if err != nil {
		return err
	}

	if !regex.MatchString(password) {
		return errors.New("")
	}

	return nil
}
