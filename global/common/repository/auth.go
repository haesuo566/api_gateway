package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/novel/auth/global/common/entity"
	util "github.com/novel/auth/global/util/sql"
)

type AuthRepository struct {
	SqlUtil util.ISqlUtil
}

var instance *AuthRepository = nil

func NewAuthRepository() *AuthRepository {
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

	// row가 없을경우 err는 리턴하지 않음
	if !rows.Next() {
		return nil, nil
	}

	user := entity.User{}
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Email, user.Credential, user.AccessToken,
			user.RefreshToken, &user.CreatedAt, &user.UpdatedAt, &user.Provider)
	}

	return &user, nil
}

// select query로 값이 있으면 update???
func (a *AuthRepository) Save(user *entity.User) (*entity.User, error) {
	var query string
	var param []any
	if user.Credential == nil { // oauth user
		query = "INSERT INTO user (name, email, oauth_access_token, oauth_refresh_token, provider) VALUES (?, ?, ?, ?, ?);"
		param = []any{user.Name, user.Email, *user.AccessToken, *user.RefreshToken, user.Provider}
	} else { // normal user
		query = "INSERT INTO user (name, email, credential, provider) VALUES (?, ?, ?, ?);"
		param = []any{user.Name, user.Email, *user.Credential, user.Provider}
	}

	// ---------------- 여기부터 transaction 으로 시작해야함 ----------------
	data, err := a.SqlUtil.QueryWithTransaction(func(tx *sql.Tx) (interface{}, error) {
		result, err := tx.Exec(query, param...)
		if err != nil {
			return nil, err
		}

		// insert가 적용되지 않음
		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		user := &entity.User{}
		row := tx.QueryRow("SELECT id, name, email, credential, oauth_access_token, oauth_refresh_token,"+
			"created_at, updated_at, provider FROM user WHERE id = ?", id).Scan(&user.Id, &user.Name, &user.Email,
			&user.Credential, &user.AccessToken, &user.RefreshToken, &user.CreatedAt, &user.UpdatedAt, &user.Provider)
		if row != nil {
			return nil, errors.New(row.Error())
		}

		return user, nil
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return data.(*entity.User), nil
}
