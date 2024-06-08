package user

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/novel/api-gateway/db"
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
	var rows *sql.Rows
	var err error
	// QueryRow를 사용하지 않고 Query를 사용한 이유는 결과가 없을떄의 에러처리가 애매함
	if tx == nil {
		rows, err = u.db.QueryContext(context.Background(), query, user.Email, user.Provider)
	} else {
		rows, err = tx.QueryContext(context.Background(), query, user.Email, user.Provider)
	}
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, &NoSearchUserError{user.Email, user.Provider}
	}

	selectedUser := User{}
	if err := rows.Scan(&selectedUser.Id, &selectedUser.Name, &selectedUser.Email, &selectedUser.Credential,
		&selectedUser.CreatedAt, &selectedUser.UpdatedAt, &selectedUser.Provider); err != nil {
		return nil, err
	}

	return &selectedUser, nil
}

// update 가능한 컬럼은 다 업데이트 하는 걸로 해야겠는데?
// user 정보를 업데이트 할때는 이 function을 사용하고
// 로그인할때만 save function을 사용함
// 이거는 개선의 여지가 많음
func (u *UserRepository) Update(user *User, tx db.ITx) (*User, error) {
	updateFunc := func(tx db.ITx) (interface{}, error) {
		selectQuery := `SELECT name, email, credential, provider FROM novel_user WHERE email = ? and provider = ?`
		row := tx.QueryRowContext(context.Background(), selectQuery, user.Email, user.Provider)
		if err := row.Err(); err != nil {
			return nil, err
		}

		originUser := &User{}
		if err := row.Scan(&originUser.Name, &originUser.Email, &originUser.Credential, &originUser.Provider); err != nil {
			return nil, err
		}

		// 시간만 업데이트 --> 로그인 기준으로만 최근 접속일을 체크함 -> 로그아웃도 체크를 해야하나?
		// if originUser.HashCode() == user.HashCode() {
		// 	updateQuery := "update novel_user set updated_at = CURRENT_TIMESTAMP() where email = ? and provider = ?"
		// 	if _, err := tx.ExecContext(context.Background(), updateQuery, originUser.Email, originUser.Provider); err != nil {
		// 		return nil, err
		// 	}

		// 	findUser, err := u.FindByEmailAndProvider(user, tx)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// 	return findUser, nil
		// }

		var buffer bytes.Buffer
		buffer.WriteString("update novel_user set updated_at = CURRENT_TIMESTAMP()")

		// 유저 이름 변경 -> 중복된 이름이 존재하는지 확인 해야함
		if originUser.Name != user.Name {
			if rows, err := tx.QueryContext(context.Background(), "select name from novel_user where name = ?", user.Name); err != nil {
				if rows.Next() {
					return nil, &DuplicateUserNameError{}
				} else {
					return nil, err
				}
			}

			buffer.WriteString(fmt.Sprintf(", name = '%s'", user.Name))
		}

		// nil pointer dereference occured
		if originUser.Credential != nil && user.Credential != nil {
			if *originUser.Credential != *user.Credential {
				buffer.WriteString(fmt.Sprintf(", credential = '%s'", *user.Credential))
			}
		}

		// 근데 이건 사용할지는 모르겠다??
		// if originUser.Email != user.Email {

		// }
		buffer.WriteString(" where email = ? and provider = ?")
		if _, err := tx.ExecContext(context.Background(), buffer.String(), user.Email, user.Provider); err != nil {
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
		// 유저가 없는 경우
		if _, err := u.FindByEmailAndProvider(user, tx); reflect.TypeOf(err) == reflect.TypeOf(&NoSearchUserError{}) {
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
		} else if err != nil { // 그냥 에러인 경우
			return nil, err
		} else { // 유저가 있는 경우
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
