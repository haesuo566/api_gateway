package repository

import (
	"database/sql"
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

	// row가 없을경우 err는 리턴하지 않음
	if !rows.NextResultSet() {
		return nil, nil
	}

	user := entity.User{}
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.Credential, &user.AccessToken,
			&user.RefreshToken, &user.CreatedAt, &user.UpdatedAt, &user.Provider)
	}

	return &user, nil
}

// 이 메서드를 select하고 있다면?? 업데이트를 때릴까??
// 근데 그렇게 하려면 transaction으로 해야하는데...
func (a *AuthRepository) Save(user *entity.User) (*entity.User, error) {
	var query string
	var param []any
	if user.Credential == "" { // oauth user
		query = "INSERT INTO user (name, email, oauth_access_token, oauth_refresh_token, provider) VALUES (?, ?, ?, ?, ?);"
		param = []any{user.Name, user.Email, user.AccessToken, user.RefreshToken, user.Provider}
	} else { // normal user
		query = "INSERT INTO user (name, email, credential, provider) VALUES (?, ?, ?, ?);"
		param = []any{user.Name, user.Email, user.Credential, user.Provider}
	}

	// ---------------- 여기부터 transaction 으로 시작해야함 ----------------
	a.SqlUtil.ExecWithTransaction(func(tx *sql.Tx) error {
		result, err := tx.Exec(query, param...)
		if err != nil {
			return err
		}

		// insert가 적용되지 않음
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}

		user := &entity.User{}
		tx.QueryRow("SELECT id, name, email, credential, oauth_access_token, oauth_refresh_token,"+
			"created_at, updated_at, provider FROM user WHERE id = ?", id).Scan(&user.Id, &user.Name, &user.Email,
			&user.Credential, &user.AccessToken, &user.RefreshToken, &user.CreatedAt, &user.UpdatedAt, &user.Provider)

		return nil
	})

	result, err := a.SqlUtil.Exec(query, param...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if _, err := result.RowsAffected(); err != nil {
		log.Println(err)
		return nil, err
	}

	findUser, err := a.FindById(user.Email)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return findUser, nil
}
