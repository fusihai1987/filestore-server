package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)


var db *sql.DB

func init(){

	db, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3308)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	if err := db.Ping();err != nil {
		fmt.Println("Failed to connect to mysql err: %s", err.Error)
		return
	}

	fmt.Println("Connect to mysql success!")
}

func GetDb() *sql.DB {
	return db
}

func ParseRows(rows *sql.Rows) []map[string]interface{}{
	columns,_ := rows.Columns()

	scanArgs := make([]interface{},len(columns))
	scanColumn := make([]interface{}, len(columns))

	for i := range scanColumn{
		scanArgs[i] = &scanColumn[i]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range scanColumn {
			record[columns[i]] = col
		}

		records = append(records,record)
	}
	return records
}

func checkErr(err error){
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
