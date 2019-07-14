package meta

import (
	"filestore-server/db/mysql"
	"database/sql"
	"fmt"
)

type FileDb struct {
	FileName sql.NullString
	FileSha1 sql.NullString
	FileSize sql.NullInt64
	FilePath sql.NullString
}

func Insert(fileMeta FileMeta) bool{
	db := mysql.GetDb()

	if db == nil {
		fmt.Println("Failed to init db class")
	}

	stmt, err := db.Prepare("Insert file_metas set file_name=?,file_sha1=?,file_size=?,file_path=?")
	defer stmt.Close()

	if err != nil {
		fmt.Println("Failed to prepare %s", err.Error())
		return false
	}

	res,err := stmt.Exec(fileMeta.FileName, fileMeta.FileSha1, fileMeta.FileSize, fileMeta.FilePath)

	if err != nil {
		fmt.Println("Failed to exec %s", err.Error())
		return false
	}

	if rowAffected, err:= res.RowsAffected(); err == nil {
		if rowAffected <= 0 {
			fmt.Println("Failed with has %s has been updated before",fileMeta.FileSha1)
			return false
		}

		return true
	}
	return false;
}



func GetFileDb(fileSha1 string)(*FileDb, error){

	fileDb := FileDb{}

	stmt, err := mysql.GetDb().Prepare("select file_sha1,file_name, file_size, file_path from file_metas where file_sha1=? and status = 0 limit 1")

	if err != nil {
		return &fileDb, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(fileSha1).Scan(&fileDb.FileSha1, &fileDb.FileName, &fileDb.FileSize, &fileDb.FilePath)


	if err != nil {
		return &fileDb, err
	}

	return &fileDb,nil
}
