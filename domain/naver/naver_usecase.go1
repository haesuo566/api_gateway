package naver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/novel/auth/global/common/dto"
	"github.com/novel/auth/global/common/entity"
	"github.com/novel/auth/global/common/repository"
	"github.com/novel/auth/global/util/jwt"
	"golang.org/x/oauth2"
)

type NaverUsecase struct {
	authRepository repository.IAuthRepository
	jwtUtil        jwt.IJwtUtil
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

var usecaseInstance *NaverUsecase = nil

func NewUsecase() *NaverUsecase {
	if usecaseInstance == nil {
		usecaseInstance = &NaverUsecase{
			authRepository: repository.NewAuthRepository(),
			jwtUtil:        jwt.New(),
		}
	}
	return usecaseInstance
}

func (n *NaverUsecase) GetUserInfo(token *oauth2.Token) (*entity.User, error) {
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

	user := &entity.User{
		Name:         userInfo.Response.Name,
		Email:        userInfo.Response.Email,
		AccessToken:  &token.AccessToken,
		RefreshToken: &token.RefreshToken,
		Provider:     entity.NAVER,
	}

	findUser, err := n.authRepository.FindById(user.Email)
	if err != nil {
		return nil, err
	} else if findUser == nil {
		if findUser, err = n.authRepository.Save(user); err != nil {
			return nil, err
		}
	}

	return findUser, nil
}

func (n *NaverUsecase) CreateUserToken(user *entity.User) (*dto.ResponseToken, error) {
	responseToken := &dto.ResponseToken{}

	if accessToken, err := n.jwtUtil.GenerateAccessToken(user.Email); err != nil {
		log.Println(err)
		return nil, err
	} else {
		responseToken.AccessToken = accessToken
	}

	if refreshToken, err := n.jwtUtil.GenerateRefreshToken(); err != nil {
		log.Println(err)
		return nil, err
	} else {
		responseToken.RefreshToken = refreshToken
	}

	return responseToken, nil
}
