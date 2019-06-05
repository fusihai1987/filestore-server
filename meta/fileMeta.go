package meta

import (
	"fmt"
)

type FileMeta struct {
	FileName string
	FileSha1 string
	FileSize int64
	FilePath string
	UpdatedAt string
}

var fileMetas map[string]FileMeta

func init(){
	fileMetas = make(map[string]FileMeta)
}
func Update(file FileMeta) {
	//fileMetas[file.FileSha1] = file;
	isSuccess := Insert(file)

	fmt.Println(isSuccess)
}

func GetFile(fileSha1 string) FileMeta{
	return fileMetas[fileSha1]
}