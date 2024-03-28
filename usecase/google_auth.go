package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/novel/auth/entity"
	"golang.org/x/oauth2"
)

type GoogleAuthUsecase struct {
}

var instance IGoogleAuthUsecase = nil

func NewGoogleAuthUsecase() IGoogleAuthUsecase {
	if instance == nil {
		instance = &GoogleAuthUsecase{}
	}
	return instance
}

// Signup Signin 구분 필요
func (g *GoogleAuthUsecase) GetUserInfo(token *oauth2.Token) (*entity.User, error) {
	// User Information Url
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken)

	// User Infomation Request
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Google user Info Data Type
	user := struct {
		Id            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Picture       string `json:"picture"`
		Hd            string `json:"hd"`
	}{}

	// marshaling error
	if err := json.Unmarshal(data, &user); err != nil {
		log.Println(err)
		return nil, err
	}

	return &entity.User{
		Id:       user.Email,
		Platform: entity.GOOGLE,
	}, nil
}
