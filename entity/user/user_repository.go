package user

import (
	"bytes"
	"context"

	"github.com/novel/auth/db"
)

type IUserRepository interface {
	FindByEmail(string) (*User, error)
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

func (u *UserRepository) FindByEmail(email string) (*User, error) {
	query := `SELECT id, name, email, credential, oauth_access_token, oauth_refresh_token, created_at, updated_at, provider 
	FROM user WHERE email = ?`
	row := u.db.QueryRowContext(context.Background(), query, email)
	if err := row.Err(); err != nil {
		return nil, err
	}

	user := User{}
	row.Scan(&user.Id, &user.Name, &user.Email, &user.Credential, &user.AccessToken,
		&user.RefreshToken, &user.CreatedAt, &user.UpdatedAt, &user.Provider)
	return &user, nil
}

// 아직 다 못함
// bytebuffer에 string 추가하는 부분 해야함
func (u *UserRepository) Update(user *User) (*User, error) {
	// 변경되는 컬럼 부분을 bytebuffer가?? 이걸로 붙여서 만들어
	// user select 해가지고 비교해서 변경 안된거만
	db.WithTx(u.db, func(tx db.ITx) (interface{}, error) {
		selectQuery := `SELECT name, oauth_access_token, oauth_refresh_token, updated_at FROM user WHERE email = ?`
		row := tx.QueryRowContext(context.Background(), selectQuery, user.Email)
		if err := row.Err(); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}

		u := User{}
		if err := row.Scan(&u.Name, &u.AccessToken, &u.RefreshToken, &u.UpdatedAt); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}

		var params []interface{} = make([]interface{}, 1, 4)
		var buffer bytes.Buffer
		buffer.WriteString("UPDATE user SET ")
		if u.Name != user.Name {
			buffer.WriteString("name = ?")
			params = append(params, u.Name)
		}
		if *u.AccessToken != *user.AccessToken {
			buffer.WriteString(", oauth_access_token = ?")
			params = append(params, *u.AccessToken)
		}
		if *u.RefreshToken != *user.RefreshToken {
			buffer.WriteString(", oauth_refresh_token = ?")
		}
		if u.UpdatedAt != user.UpdatedAt {
			buffer.WriteString(", oauth_refresh_token = ?")
		}

		buffer.WriteString("WEHRE email = ?")
		params = append(params, u.Email)
		updateQuery := buffer.String()
		tx.ExecContext(context.Background(), updateQuery, params...)
		return nil, nil
	})
	return nil, nil
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
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		selectQuery := `SELECT id, name, email, credential, oauth_access_token, oauth_refresh_token, provider FROM user WHERE id = ?`
		row := tx.QueryRowContext(context.Background(), selectQuery, id)
		if err := row.Err(); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}

		iUser := User{}
		if err := row.Scan(&iUser.Id, &iUser.Name, &iUser.Email, &iUser.Credential, &iUser.AccessToken, &iUser.RefreshToken, &iUser.Provider); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}

		return &iUser, nil
	})

	if err != nil {
		return nil, err
	}

	return newUser.(*User), nil
}
