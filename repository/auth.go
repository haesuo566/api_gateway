package repository

import (
	"github.com/novel/auth/entity"
	"github.com/novel/auth/util/sql"
)

type AuthRepository struct {
	sql sql.ISqlUtil
}

var instance IAuthRepository = nil

func NewAuthRepository() IAuthRepository {
	if instance == nil {
		instance = &AuthRepository{
			sql: sql.New(),
		}
	}
	return instance
}

func (a *AuthRepository) FindById(id string) (*entity.User, error) {

	return nil, nil
}
