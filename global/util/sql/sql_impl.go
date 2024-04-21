package sql

import "database/sql"

type ISqlUtil interface {
	Query(query string, param ...interface{}) (*sql.Rows, error)
	Exec(query string, param ...interface{}) (sql.Result, error)
	QueryWithTransaction(transaction transactionWithResult) (interface{}, error)
	ExecWithTransaction(transaction transaction) error
}

type transactionWithResult func(*sql.Tx) (interface{}, error)
type transaction func(*sql.Tx) error
