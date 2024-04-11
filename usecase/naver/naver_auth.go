package naver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/novel/auth/entity"
	"github.com/novel/auth/usecase"
	"github.com/novel/auth/util/jwt"
	"golang.org/x/oauth2"
)

type naverAuthUsecase struct {
}

type NaverUserInfo struct {
	ResultCode string `json:"resultCode"`
	Message    string `json:"message"`
	Response   struct {
		Id    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"response"`
}

var instance usecase.IAuthUsecase = nil

func New() usecase.IAuthUsecase {
	if instance == nil {
		instance = &naverAuthUsecase{}
	}
	return instance
}

func (n *naverAuthUsecase) GetUserInfo(token *oauth2.Token) (*entity.User, error) {
	accessToken := token.AccessToken
	url := "https://openapi.naver.com/v1/nid/me"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	userInfo := &NaverUserInfo{}
	if err := json.Unmarshal(data, userInfo); err != nil {
		log.Println(err)
		return nil, err
	}

	fmt.Println(userInfo)
	return nil, nil
}

func (n *naverAuthUsecase) CreateUserToken(user *entity.User) (*jwt.ResponseToken, error) {
	return nil, nil
}
