package google

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/novel/auth/entity/user"
	"golang.org/x/oauth2"
)

type IGoogleUsecase interface {
	GetUserInfo(*oauth2.Token) (*user.User, error)
}

type GoogleUsecase struct {
	userRepository user.IUserRepository
}

type googleUserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

var usecaseInstance IGoogleUsecase = nil

func NewUsecase(userRepository user.IUserRepository) IGoogleUsecase {
	if usecaseInstance == nil {
		usecaseInstance = &GoogleUsecase{
			userRepository: userRepository,
		}
	}
	return usecaseInstance
}

// 로그인하면 updated_at column 업데이트 하는 부분 추가
func (g *GoogleUsecase) GetUserInfo(token *oauth2.Token) (*user.User, error) {
	// User Information Url
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken)

	// User Infomation Request
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	googleUserInfo := googleUserInfo{}
	if err := json.Unmarshal(data, &googleUserInfo); err != nil {
		return nil, err
	}

	user := &user.User{
		Name:     googleUserInfo.Name,
		Email:    googleUserInfo.Email,
		Provider: user.GOOGLE,
	}

	saveUser, err := g.userRepository.Save(user)
	if err != nil {
		return nil, err
	}

	return saveUser, nil
}
