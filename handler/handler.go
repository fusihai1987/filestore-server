package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET"{
		uploadHtml, err := ioutil.ReadFile("./static/views/upload.html")

		if err != nil {
			io.WriteString(w, "File store server internal error!")
			return
		}

		io.WriteString(w,string(uploadHtml))
	}else if r.Method == "POST"{
		file,header, err :=  r.FormFile("file")

		if err != nil{
			fmt.Println("Failed to get file:s%", err.Error())
			return
		}
		defer file.Close()

		localFile, err := os.Create("./tmp/"  + header.Filename)

		if err != nil{
			fmt.Println("Failed to create file:%s\n", err.Error())
			return
		}
		defer localFile.Close()

		_, err = io.Copy(localFile, file)

		if err != nil{
			fmt.Println("Failed to save file :%s\n", err.Error())
			return
		}

		io.WriteString(w,fmt.Sprintf("{\"code\":200, \"msg\":success,\"url\":/tmp/%s}", header.Filename))
	}

}

func SucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Success")
}


