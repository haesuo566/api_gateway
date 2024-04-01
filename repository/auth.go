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

func (a *AuthRepository) FindById(email string) (*entity.User, error) {
	query := "SELECT id, name, email, credential, oauth_access_token, oauth_refresh_token," +
		"created_at, updated_at, provider FROM user WHERE email = ?"
	rows, err := a.SqlUtil.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := entity.User{}
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.Credential, &user.OauthAccessToken,
			&user.OauthRefreshToken, &user.CreatedAt, &user.UpdatedAt, &user.Provider)
	}

	return &user, nil
}
