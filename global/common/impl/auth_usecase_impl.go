package impl

import (
	"github.com/novel/auth/global/common/dto"
	"github.com/novel/auth/global/common/entity"
	"golang.org/x/oauth2"
)

type AuthImpl interface {
	GetUserInfo(token *oauth2.Token) (*entity.User, error)
	CreateUserToken(user *entity.User) (*dto.ResponseToken, error)
}
