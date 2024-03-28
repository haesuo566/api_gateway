package usecase

import "golang.org/x/oauth2"

type IGoogleAuthUsecase interface {
	GetUserInfo(token *oauth2.Token) (interface{}, error)
}
