package db

import (
	db "FILESTORE-SERVER/db/mysql"
	"FILESTORE-SERVER/db/mysql"
	"FILESTORE-SERVER/utils"
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
		return false
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

func Signin(username string, passwd string) bool{
	stmt, err := db.GetDb().Prepare("select * from user where user_name=? limit 1")

	if err != nil {
		fmt.Println("Db parepare err :", err.Error())
		return false
	}

	defer stmt.Close()

	rows,err:= stmt.Query(username)

	if err != nil {
		fmt.Println("Query failed:", err.Error())
		return false
	}else if rows == nil {
		fmt.Println("Rows empty:")
		return false
	}

	pRows := db.ParseRows(rows)

	fmt.Println(string(pRows[0]["user_pwd"].([]byte)))
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == utils.MD5([]byte(passwd)){
		return true
	}

	return false
}

func GetUserInfo(username string) (*User,error){
	user := User{}

	stmt, err := mysql.GetDb().Prepare("select user_name, created_at from user where user_name =? limit 1")

	if err != nil {
		return &user, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.CreatedAt)


	if err != nil {
		return &user, err
	}

	return &user,nil
}

