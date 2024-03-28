package usecase

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type GoogleAuthUsecase struct {
}

func NewGoogleAuthUsecase() IGoogleAuthUsecase {
	return &GoogleAuthUsecase{}
}

func (g *GoogleAuthUsecase) GetUserInfo(token *oauth2.Token) (interface{}, error) {
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

	log.Println(string(data))

	// token.AccessToken
	return nil, nil
}
