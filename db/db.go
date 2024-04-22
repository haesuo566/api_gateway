package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

type IDatabase interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	Begin() (*sql.Tx, error)
}

var instance IDatabase = nil

func New() IDatabase {
	if instance == nil {
		format := "%s:%s@tcp(%s:%s)/auth?parseTime=True"
		url := fmt.Sprintf(format, os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
		db, err := sql.Open("mysql", url)
		if err != nil {
			panic(err)
		}

		db.SetMaxIdleConns(10) // connection pool 갯수를 의미(최대 대기 가능한 connection 갯수)
		db.SetMaxOpenConns(10) // 한번에 사용 가능한 connection의 갯수를 의미

		instance = db
	}

	return instance
}
