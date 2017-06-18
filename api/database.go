package api

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var conn *sql.DB = nil

func InitDatabaseConnection() {
	db, err := sql.Open("mysql", MYSQL_CONN_STR)
	checkErr(err)
	err = db.Ping()
	checkErr(err)

	conn = db
}
