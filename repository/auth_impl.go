package repository

import "github.com/novel/auth/entity"

type IAuthRepository interface {
	FindById(id string) (*entity.User, error)
	Save(user *entity.User) (*entity.User, error)
}
