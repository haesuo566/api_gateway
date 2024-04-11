package usecase

import (
	"github.com/novel/auth/entity"
	"github.com/novel/auth/util/jwt"
	"golang.org/x/oauth2"
)

type IAuthUsecase interface {
	GetUserInfo(token *oauth2.Token) (*entity.User, error)
	CreateUserToken(user *entity.User) (*jwt.ResponseToken, error)
}
