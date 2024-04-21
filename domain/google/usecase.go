package google

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/novel/auth/global/common/dto"
	"github.com/novel/auth/global/common/entity"
	"github.com/novel/auth/global/common/repository"
	"github.com/novel/auth/global/util/jwt"
	"golang.org/x/oauth2"
)

type GoogleAuthUsecase struct {
	authRepository repository.IAuthRepository
	jwtUtil        jwt.IJwtUtil
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

var usecaseInstance *GoogleAuthUsecase = nil

func NewUsecase() *GoogleAuthUsecase {
	if usecaseInstance == nil {
		usecaseInstance = &GoogleAuthUsecase{
			authRepository: repository.NewAuthRepository(),
			jwtUtil:        jwt.New(),
		}
	}
	return usecaseInstance
}

// Signup Signin 구분 필요 -> 처리 완룐
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

	gUser := googleUserInfo{}
	if err := json.Unmarshal(data, &gUser); err != nil {
		log.Println(err)
		return nil, err
	}

	user := &entity.User{
		Name:         gUser.Name,
		Email:        gUser.Email,
		AccessToken:  &token.AccessToken,
		RefreshToken: &token.RefreshToken,
		Provider:     entity.GOOGLE,
	}

	findUser, err := g.authRepository.FindById(user.Email)
	if err != nil {
		return nil, err
	} else if findUser == nil {
		if findUser, err = g.authRepository.Save(user); err != nil {
			return nil, err
		}
	}

	return findUser, nil
}

func (g *GoogleAuthUsecase) CreateUserToken(user *entity.User) (*dto.ResponseToken, error) {
	responseToken := &dto.ResponseToken{}

	if accessToken, err := g.jwtUtil.GenerateAccessToken(user.Email); err != nil {
		log.Println(err)
		return nil, err
	} else {
		responseToken.AccessToken = accessToken
	}

	if refreshToken, err := g.jwtUtil.GenerateRefreshToken(); err != nil {
		log.Println(err)
		return nil, err
	} else {
		responseToken.RefreshToken = refreshToken
	}

	return responseToken, nil
}
