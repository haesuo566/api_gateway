package user

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"github.com/novel/auth/db"
)

type IUserRepository interface {
	FindByEmail(string, db.ITx) (*User, error)
	Update(*User) (*User, error)
	Save(*User) (*User, error)
}

type UserRepository struct {
	db db.IDatabase
}

var instance IUserRepository = nil

func NewRepository(db db.IDatabase) IUserRepository {
	if instance == nil {
		instance = &UserRepository{
			db: db,
		}
	}
	return instance
}

func (u *UserRepository) FindByEmail(email string, tx db.ITx) (*User, error) {
	query := `SELECT id, name, email, credential, oauth_access_token, oauth_refresh_token, created_at, updated_at, provider 
	FROM user WHERE email = ?`
	var row *sql.Row
	if tx == nil {
		row = u.db.QueryRowContext(context.Background(), query, email)
	} else {
		row = tx.QueryRowContext(context.Background(), query, email)
	}

	if err := row.Err(); err != nil { // 이거 에러처리가 애매하다
		return nil, err
	}

	user := User{}
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Credential, &user.AccessToken,
		&user.RefreshToken, &user.CreatedAt, &user.UpdatedAt, &user.Provider); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) Update(user *User) (*User, error) {
	_, err := db.WithTx(u.db, func(tx db.ITx) (interface{}, error) {
		selectQuery := `SELECT name, oauth_access_token, oauth_refresh_token, updated_at FROM user WHERE email = ?`
		row := tx.QueryRowContext(context.Background(), selectQuery, user.Email)
		if err := row.Err(); err != nil {
			return nil, err
		}

		origUser := User{}
		if err := row.Scan(&origUser.Name, &origUser.AccessToken, &origUser.RefreshToken, &origUser.UpdatedAt); err != nil {
			return nil, err
		}

		// var columnName []string = []string{"name", "updated_at"}
		// var originUser []interface{} = []interface{}{origUser.Name, origUser.UpdatedAt}
		// var newUser []interface{} = []interface{}{user.Name, user.UpdatedAt}
		var params []interface{} = make([]interface{}, 0, 4)
		var buffer bytes.Buffer
		buffer.WriteString("UPDATE user SET ")
		var b bool = false

		combineQuery := func(columnName string, first interface{}, second interface{}, param interface{}) {
			if first != second {
				if b {
					buffer.WriteString(", ")
				}

				buffer.WriteString(fmt.Sprintf("%s = ?", columnName))
				params = append(params, param)
				b = true
			}
		}

		combineQuery("name", origUser.Name, user.Name, user.Name)
		combineQuery("updated_at", origUser.UpdatedAt, user.UpdatedAt, "NOW()")

		buffer.WriteString(" WHERE email = ?")
		updateQuery := buffer.String()
		params = append(params, user.Email)
		_, err := tx.ExecContext(context.Background(), updateQuery, params...)
		if err != nil {
			return nil, err
		}

		findUser, err := u.FindByEmail(user.Email, tx)
		if err != nil {
			return nil, err
		}

		return findUser, nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) Save(user *User) (*User, error) {
	var query string
	var params []interface{}
	if user.Credential != nil {
		query = "INSERT INTO user(name, email, credential, provider) VALUES(?, ?, ?, ?)"
		params = []interface{}{user.Name, user.Email, *user.Credential, user.Provider}
	} else {
		query = "INSERT INTO user(name, email, oauth_access_token, oauth_refresh_token, provider) VALUES(?, ?, ?, ?, ?, ?)"
		params = []interface{}{user.Name, user.Email, *user.AccessToken, *user.RefreshToken, user.Provider}
	}

	newUser, err := db.WithTx(u.db, func(tx db.ITx) (interface{}, error) {
		result, err := tx.ExecContext(context.Background(), query, params...)
		if err != nil {
			return nil, err
		}

		if _, err := result.LastInsertId(); err != nil {
			return nil, err
		}

		newUser, err := u.FindByEmail(user.Email, tx)
		if err != nil {
			return nil, err
		}

		return newUser, nil
	})

	if err != nil {
		return nil, err
	}

	return newUser.(*User), nil
}
