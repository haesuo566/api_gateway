package sql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/novel/auth/config"
)

type SqlUtil struct {
	db *sql.DB
}

var instance ISqlUtil = nil

func New() ISqlUtil {
	if instance == nil {
		host := config.Getenv("DB_HOST")
		port := config.Getenv("DB_PORT")
		user := config.Getenv("DB_USER")
		pass := config.Getenv("DB_PASS")

		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/auth", user, pass, host, port))
		if err != nil {
			log.Println(err)
			return nil
		}

		db.SetConnMaxIdleTime(time.Minute * 1) // connection pool timeout
		db.SetMaxIdleConns(10)                 // connection pool
		db.SetMaxOpenConns(100)                // active connection

		instance = &SqlUtil{
			db: db,
		}
	}
	return instance
}

func (s *SqlUtil) Query(query string, param ...interface{}) (*sql.Rows, error) {
	// s.db.QueryRow
	rows, err := s.db.Query(query, param...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// if !rows.NextResultSet() {
	// 	return nil, nil
	// }
	log.Println(query)

	return rows, nil
}

func (s *SqlUtil) Exec(query string, param ...interface{}) (sql.Result, error) {
	result, err := s.db.Exec(query, param...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(query)

	return result, nil
}

// query, exec를 variadic argument로 받아할 듯
func (s *SqlUtil) ExecWithTransaction(transaction Transaction) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	if err := transaction(tx); err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}

	return nil
}
