package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/novel/api-gateway/config"
)

type ITx interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type IDatabase interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	Begin() (*sql.Tx, error)
}

var (
	instance IDatabase = nil
	// originInstance *sql.DB   = nil
)

func New() IDatabase {
	if instance == nil {
		config.LoadEnv()
		format := "%s:%s@tcp(%s:%s)/auth?parseTime=True"
		url := fmt.Sprintf(format, os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
		db, err := sql.Open("mysql", url)
		if err != nil {
			panic(err)
		}

		db.SetMaxIdleConns(10) // connection pool 갯수를 의미(최대 대기 가능한 connection 갯수)
		db.SetMaxOpenConns(10) // 한번에 사용 가능한 connection의 갯수를 의미

		instance = db
		// originInstance = db
	}

	return instance
}

func WithTx(db IDatabase, queries func(tx ITx) (interface{}, error)) (interface{}, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	data, err := queries(tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	return data, nil
}
