package repository

import "github.com/novel/auth/global/common/entity"

type IAuthRepository interface {
	FindById(id string) (*entity.User, error)
	Save(user *entity.User) (*entity.User, error)
}
