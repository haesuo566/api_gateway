package user

import (
	"context"

	"github.com/novel/auth/db"
)

type IUserRepository interface {
	FindByEmail(string) (*User, error)
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

func (u *UserRepository) Save(user *User) (*User, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	return nil, nil
}
