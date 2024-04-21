package google

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/novel/auth/entity"
	"golang.org/x/oauth2"
)

type IGoogleUsecase interface {
	GetUserInfo(token *oauth2.Token) (*entity.User, error)
}

type GoogleUsecase struct {
	repository IGoogleRepository
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

func NewUsecase(repository IGoogleRepository) IGoogleUsecase {
	if usecaseInstance == nil {
		usecaseInstance = &GoogleUsecase{
			repository: repository,
		}
	}
	return usecaseInstance
}

func (g *GoogleUsecase) GetUserInfo(token *oauth2.Token) (*entity.User, error) {
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

	user := &entity.User{
		Name:         googleUserInfo.Name,
		Email:        googleUserInfo.Email,
		AccessToken:  &token.AccessToken,
		RefreshToken: &token.RefreshToken,
		Provider:     entity.GOOGLE,
	}

	findUser, err := g.repository.FindById(user.Email)
	if err != nil {
		return nil, err
	} else if findUser == nil {
		if findUser, err = g.repository.Save(user); err != nil {
			return nil, err
		}
	}

	return findUser, nil
}
