package usecase

import (
	"github.com/novel/auth/entity"
	"golang.org/x/oauth2"
)

type IGoogleAuthUsecase interface {
	GetUserInfo(token *oauth2.Token) (*entity.User, error)
}
