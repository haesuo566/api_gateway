package google

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/novel/auth/entity"
	"github.com/novel/auth/repository"
	"github.com/novel/auth/usecase"
	"github.com/novel/auth/util/jwt"
	"github.com/novel/auth/util/sql"
	"golang.org/x/oauth2"
)

type GoogleAuthUsecase struct {
	googleAuthRepository repository.IAuthRepository
	jwtUtil              jwt.IJwtUtil
	sqlUtil              sql.ISqlUtil
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

var instance usecase.IAuthUsecase = nil

func New() usecase.IAuthUsecase {
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

func (g *GoogleAuthUsecase) CreateUserToken(user *entity.User) (*jwt.ResponseToken, error) {
	responseToken := &jwt.ResponseToken{}

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
