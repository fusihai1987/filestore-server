package db

import (
	db "FILESTORE-SERVER/db/mysql"
	"fmt"
)

type User struct{
	Username string
	CreatedAt string
}

func SignUp(username string, passwd string) bool {
	stmt, err := db.GetDb().Prepare("INSERT  user  set user_name=?,user_pwd=?")

	if err != nil {
		fmt.Println("Db parpare err: %s", err.Error())
	}
	defer stmt.Close()


	res, err := stmt.Exec(username, passwd)

	if err != nil {
		fmt.Println("Failed to exec : %s", err.Error())
	}

	if rowaffected, err := res.RowsAffected(); err == nil && rowaffected > 0 {
		return true
	}

	return false

}