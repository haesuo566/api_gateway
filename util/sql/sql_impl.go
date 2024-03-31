package sql

import "database/sql"

type ISqlUtil interface {
	Query(query string, param ...interface{}) (*sql.Rows, error)
	Exec(query string, param ...interface{}) (*sql.Result, error)
}
