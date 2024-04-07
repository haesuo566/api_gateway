package usecase

import (
	"github.com/novel/auth/entity"
	"github.com/novel/auth/util/jwt"
	"golang.org/x/oauth2"
)

type IGoogleAuthUsecase interface {
	GetUserInfo(token *oauth2.Token) (*entity.User, error)
	CreateUserToken(user *entity.User) (*jwt.ResposneToken, error)
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
