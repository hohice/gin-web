package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hohice/gin-web/db/mapper"
)

var db *sql.DB

func InitDB() (err error) {
	if db, err = sql.Open("mysql", "user:password@/dbname"); err != nil {
		return err
	}
	return
}

func Close() {
	db.Close()
}

func GetUserMapper() *mapper.UserMapper {
	return &mapper.UserMapper{
		DB: db,
	}
}
