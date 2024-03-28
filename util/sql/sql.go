package sql

import (
	"database/sql"
	"fmt"
	"log"

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

		// db.SetConnMaxIdleTime()
		// db.SetMaxIdleConns()
		// db.SetMaxOpenConns()

		instance = &SqlUtil{
			db: db,
		}
	}
	return instance
}

func (s *SqlUtil) Query() {
	// sql.
}

func (s *SqlUtil) Exec() {

}
