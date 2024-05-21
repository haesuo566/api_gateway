package kakao

import (
	"fmt"
	"io"
	"net/http"

	"github.com/novel/api-gateway/entity/user"
	"golang.org/x/oauth2"
)

type IKakaoUsecase interface {
	GetUserInfo(*oauth2.Token) (*user.User, error)
}

type KakaoUsecase struct {
	userRepository user.IUserRepository
}

var usecaseInstance IKakaoUsecase = nil

func NewUsecase(userRepository user.IUserRepository) IKakaoUsecase {
	if usecaseInstance == nil {
		usecaseInstance = &KakaoUsecase{
			userRepository: userRepository,
		}
	}
	return usecaseInstance
}

// 로그인하면 updated_at column 업데이트 하는 부분 추가
func (g *KakaoUsecase) GetUserInfo(token *oauth2.Token) (*user.User, error) {
	// User Information Url
	request, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	request.Header.Set("Content-Type", "application/json;charset=utf-8")

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

	fmt.Println(string(data))

	// naverUserInfo := NaverUserInfo{}
	// if err := json.Unmarshal(data, &naverUserInfo); err != nil {
	// 	return nil, err
	// }

	// user := &user.User{
	// 	Name:     naverUserInfo.Response.Name,
	// 	Email:    naverUserInfo.Response.Email,
	// 	Provider: 0,
	// }

	// saveUser, err := g.userRepository.Save(user)
	// if err != nil {
	// 	return nil, err
	// }

	// return saveUser, nil
	return nil, nil
}
