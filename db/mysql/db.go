package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)


var db *sql.DB

func init(){

	db, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3309)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	if err := db.Ping();err != nil {
		fmt.Println("Failed to connect to mysql err: %s", err.Error)
	}

	fmt.Println("Connect to mysql success!");
}

func GetDb() *sql.DB {
	return db;
}

