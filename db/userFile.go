package db

import (
	"filestore-server/db/mysql"
	"fmt"
)

type UserFile struct{
	UserName string
	FileSha1 string
	FileName string
	FileSize int64
	UploadedAt string
	UpdatedAt string
}

func Insert(userFile UserFile) bool {
	db := mysql.GetDb()

	stmt,err := db.Prepare("INSERT tbl_user_file SET user_name=?,file_sha1 = ?, file_name=?,file_size=?,uploaded_at=?,updated_at=?")

	if err != nil {
		fmt.Printf("Stmt err:%s", err.Error())
		return false
	}

	res ,err := stmt.Exec(userFile.UserName,userFile.FileSha1,userFile.FileName,userFile.FileSize,userFile.UploadedAt,userFile.UpdatedAt)

	if err != nil {
		fmt.Printf("Exec err:%s", err.Error())
		return false
	}

	if rowsAffected, err := res.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}
	return false
}

func QueryUserFile(username string, limit int) ([]UserFile,error){
	stmt, err := mysql.GetDb().Prepare("select user_name,file_name,file_sha1,file_size,uploaded_at from tbl_user_file where user_name =? limit ?")

	if err != nil {
		fmt.Println("QueryUserFile stmt prepare err", err.Error())
		return nil,err
	}

	rows,err := stmt.Query(username, limit)

	if err != nil {
		fmt.Println("QueryUserFile query err", err.Error())
		return nil,err
	}

	var userFiles []UserFile

	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.UserName, &ufile.FileName, &ufile.FileSha1,&ufile.FileSize,&ufile.UpdatedAt)
		if err != nil {
			fmt.Println("QueryUserFile next err", err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}

	return userFiles,nil
}