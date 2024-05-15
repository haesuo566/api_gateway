package user

import (
	"bytes"
	"context"
	"database/sql"

	"github.com/novel/auth/db"
)

type IUserRepository interface {
	FindByEmailAndProvider(*User, db.ITx) (*User, error)
	Update(*User, db.ITx) (*User, error)
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

func (u *UserRepository) FindByEmailAndProvider(user *User, tx db.ITx) (*User, error) {
	query := `SELECT id, name, email, credential, created_at, updated_at, provider FROM novel_user WHERE email = ? and provider = ?`
	var row *sql.Row
	if tx == nil {
		row = u.db.QueryRowContext(context.Background(), query, user.Email, user.Provider)
	} else {
		row = tx.QueryRowContext(context.Background(), query, user.Email, user.Provider)
	}

	if err := row.Err(); err != nil { // 이거 에러처리가 애매하다
		return nil, err
	}

	selectedUser := User{}
	if err := row.Scan(&selectedUser.Id, &selectedUser.Name, &selectedUser.Email, &selectedUser.Credential,
		&selectedUser.CreatedAt, &selectedUser.UpdatedAt, &selectedUser.Provider); err != nil {
		return nil, err
	}

	return &selectedUser, nil
}

func (u *UserRepository) Update(user *User, tx db.ITx) (*User, error) {
	updateFunc := func(tx db.ITx) (interface{}, error) {
		selectQuery := `SELECT name, updated_at FROM novel_user WHERE email = ?`
		row := tx.QueryRowContext(context.Background(), selectQuery, user.Email)
		if err := row.Err(); err != nil {
			return nil, err
		}

		origUser := User{}
		if err := row.Scan(&origUser.Name, &origUser.UpdatedAt); err != nil {
			return nil, err
		}

		var params []interface{} = make([]interface{}, 0, 3)
		var buffer bytes.Buffer
		buffer.WriteString("UPDATE novel_user SET updated_at = CURRENT_TIMESTAMP")

		if origUser.Name != user.Name {
			buffer.WriteString(", name = ?")
			params = append(params, user.Name)
		}

		buffer.WriteString(" WHERE email = ?")
		updateQuery := buffer.String()
		params = append(params, user.Email)
		_, err := tx.ExecContext(context.Background(), updateQuery, params...)
		if err != nil {
			return nil, err
		}

		findUser, err := u.FindByEmailAndProvider(user, tx)
		if err != nil {
			return nil, err
		}

		return findUser, nil
	}

	var userInfo interface{}
	var err error
	if tx != nil {
		userInfo, err = updateFunc(tx)
	} else {
		userInfo, err = db.WithTx(u.db, updateFunc)
	}

	if err != nil {
		return nil, err
	}

	return userInfo.(*User), nil
}

func (u *UserRepository) Save(user *User) (*User, error) {
	var query string
	var params []interface{}
	if user.Credential != nil {
		query = "INSERT INTO novel_user(name, email, credential, provider) VALUES(?, ?, ?, ?)"
		params = []interface{}{user.Name, user.Email, *user.Credential, user.Provider}
	} else {
		query = "INSERT INTO novel_user(name, email, provider) VALUES(?, ?, ?)"
		params = []interface{}{user.Name, user.Email, user.Provider}
	}

	newUser, err := db.WithTx(u.db, func(tx db.ITx) (interface{}, error) {
		var userInfo *User
		if _, err := u.FindByEmailAndProvider(user, tx); err != nil {
			result, err := tx.ExecContext(context.Background(), query, params...)
			if err != nil {
				return nil, err
			}

			if _, err := result.LastInsertId(); err != nil {
				return nil, err
			}

			if userInfo, err = u.FindByEmailAndProvider(user, tx); err != nil {
				return nil, err
			}
		} else {
			if userInfo, err = u.Update(user, tx); err != nil {
				return nil, err
			}
		}

		return userInfo, nil
	})

	if err != nil {
		return nil, err
	}

	return newUser.(*User), nil
}
