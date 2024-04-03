package repository

import (
	"log"

	"github.com/novel/auth/entity"
	util "github.com/novel/auth/util/sql"
)

type AuthRepository struct {
	SqlUtil util.ISqlUtil
}

var instance IAuthRepository = nil

func NewAuthRepository() IAuthRepository {
	if instance == nil {
		instance = &AuthRepository{
			SqlUtil: util.New(),
		}
	}
	return instance
}

func (a *AuthRepository) FindById(email string) (*entity.User, error) {
	query := "SELECT id, name, email, credential, oauth_access_token, oauth_refresh_token," +
		"created_at, updated_at, provider FROM user WHERE email = ?;"
	rows, err := a.SqlUtil.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := entity.User{}
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.Credential, &user.AccessToken,
			&user.RefreshToken, &user.CreatedAt, &user.UpdatedAt, &user.Provider)
	}

	return &user, nil
}

// 이 메서드를 select하고 있다면?? 업데이트를 때릴까??
func (a *AuthRepository) Save(user *entity.User) error {
	var query string
	var param []any
	if user.Credential == "" { // oauth user
		query = "INSERT INTO user (name, email, oauth_access_token, oauth_refresh_token, provider) VALUES (?, ?, ?, ?, ?);"
		param = []any{user.Name, user.Email, user.AccessToken, user.RefreshToken, user.Provider}
	} else { // normal user
		query = "INSERT INTO user (name, email, credential, provider) VALUES (?, ?, ?, ?);"
		param = []any{user.Name, user.Email, user.Credential, user.Provider}
	}

	result, err := a.SqlUtil.Exec(query, param...)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err := result.RowsAffected(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
