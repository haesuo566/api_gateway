package naver

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/novel/api-gateway/entity/user"
	"golang.org/x/oauth2"
)

type iNaverUsecase interface {
	getUserInfo(*oauth2.Token) (*user.User, error)
}

type naverUsecase struct {
	userRepository user.IUserRepository
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

var usecaseInstance iNaverUsecase = nil

func newUsecase(userRepository user.IUserRepository) iNaverUsecase {
	if usecaseInstance == nil {
		usecaseInstance = &naverUsecase{
			userRepository: userRepository,
		}
	}
	return usecaseInstance
}

// 로그인하면 updated_at column 업데이트 하는 부분 추가
func (g *naverUsecase) getUserInfo(token *oauth2.Token) (*user.User, error) {
	// User Information Url
	request, err := http.NewRequest("GET", "https://openapi.naver.com/v1/nid/me", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	naverUserInfo := NaverUserInfo{}
	if err := json.Unmarshal(data, &naverUserInfo); err != nil {
		return nil, err
	}

	user := &user.User{
		Name:     naverUserInfo.Response.Name,
		Email:    naverUserInfo.Response.Email,
		Provider: user.NAVER,
	}

	saveUser, err := g.userRepository.Save(user)
	if err != nil {
		return nil, err
	}

	return saveUser, nil
}
