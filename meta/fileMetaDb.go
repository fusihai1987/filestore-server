package meta

import (
	"FILESTORE-SERVER/db/mysql"
	"fmt"
)

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
			fmt.Println("File with has %s has been updated before",fileMeta.FileSha1)
			return false
		}

		return true
	}
	return false;
}