package repository

import (
	"github.com/novel/auth/entity"
	"github.com/novel/auth/util/sql"
)

type AuthRepository struct {
	SqlUtil sql.ISqlUtil
}

var instance IAuthRepository = nil

func NewAuthRepository() IAuthRepository {
	if instance == nil {
		instance = &AuthRepository{
			SqlUtil: sql.New(),
		}
	}
	return instance
}

func (a *AuthRepository) FindById(id string) (*entity.User, error) {
	rows, err := a.SqlUtil.Query("SELECT * FROM user WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := &entity.User{}
	for rows.Next() {
		// rows.Scan(user.Id, user.Platform, user.A)
	}

	return user, nil
}
