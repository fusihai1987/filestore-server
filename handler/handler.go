package handler

import (
	"FILESTORE-SERVER/meta"
	"FILESTORE-SERVER/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET"{
		uploadHtml, err := ioutil.ReadFile("./static/views/upload.html")

		if err != nil {
			io.WriteString(w, "File store server internal error!")
			return
		}

		io.WriteString(w, string(uploadHtml))
	}else if r.Method == "POST"{
		file,header, err :=  r.FormFile("file")

		if err != nil{
			fmt.Println("Failed to get file:s%", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			FilePath: "./tmp/" + header.Filename,
			UpdatedAt: time.Now().Format("2006-1-2 15:04:05"),
		}

		localFile, err := os.Create(fileMeta.FilePath)

		if err != nil{
			fmt.Println("Failed to create file:%s\n", err.Error())
			return
		}
		defer localFile.Close()

		fileMeta.FileSize, err = io.Copy(localFile, file)

		if err != nil{
			fmt.Println("Failed to save file :%s\n", err.Error())
			return
		}
		localFile.Seek(0,0)
		fileMeta.FileSha1 = utils.FileSha1(localFile)
		meta.Update(fileMeta)

		fileMetaJson , err := json.Marshal(fileMeta)

		if err != nil {
			fmt.Printf("Failed to json marsha1 err:%s", err.Error())
		}
		io.WriteString(w,fmt.Sprintf("{\"code\":200, \"msg\":success,\"file\":%s}", string(fileMetaJson)))
	}

}

func QueryFile(w http.ResponseWriter, r *http.Request){
	fileSha1 := r.URL.Query()
	//fileMeta := meta.GetFile(fileSha1["filesha1"][0])
	file,err := meta.GetFileDb(fileSha1["filesha1"][0])
	fmt.Println("sha1" + fileSha1["filesha1"][0])
	if err != nil {
		fmt.Printf("Failed to get file meta:%s", err.Error())
	}

	fileMeta := meta.FileMeta{
		FileName: file.FileName.String,
		FileSha1: file.FileSha1.String,
		FileSize: file.FileSize.Int64,
		FilePath: file.FilePath.String,
	}

	fileMetaJson, err := json.Marshal(fileMeta)


	if err != nil {
		fmt.Printf("Failed to json marsha1 err:%s", err.Error())
	}
	w.Write(fileMetaJson)
}

func SucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Success")
}


