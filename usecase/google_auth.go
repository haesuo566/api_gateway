package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/novel/auth/entity"
	"github.com/novel/auth/repository"
	"github.com/novel/auth/util/jwt"
	"github.com/novel/auth/util/sql"
	"golang.org/x/oauth2"
)

type GoogleAuthUsecase struct {
	googleAuthRepository repository.IAuthRepository
	jwtUtil              jwt.IJwtUtil
	sqlUtil              sql.ISqlUtil
}

var instance IGoogleAuthUsecase = nil

func NewGoogleAuthUsecase() IGoogleAuthUsecase {
	if instance == nil {
		instance = &GoogleAuthUsecase{
			googleAuthRepository: repository.NewAuthRepository(),
			jwtUtil:              jwt.New(),
			sqlUtil:              sql.New(),
		}
	}
	return instance
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

	findUser, err := g.googleAuthRepository.Save(user)
	if err != nil {
		return nil, err
	}

	return findUser, nil
}

func (g *GoogleAuthUsecase) CreateUserToken(user *entity.User) (*jwt.ResposneToken, error) {
	responseToken := &jwt.ResposneToken{}

	if accessToken, err := jwt.New().GenerateAccessToken(user.Email); err != nil {
		log.Println(err)
		return nil, err
	} else {
		responseToken.AccessToken = accessToken
	}

	if refreshToken, err := jwt.New().GenerateRefreshToken(); err != nil {
		log.Println(err)
		return nil, err
	} else {
		responseToken.RefreshToken = refreshToken
	}

	return responseToken, nil
}
