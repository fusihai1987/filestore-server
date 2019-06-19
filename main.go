package main

import (
	"FILESTORE-SERVER/handler"
	"fmt"
	"net/http"
	"FILESTORE-SERVER/db/mysql"
)

func main(){
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/success", handler.SucHandler)
	http.HandleFunc("/file/query",handler.QueryFile)
	http.HandleFunc("/user/signup",  handler.SignUpHandler)
	mysql.GetDb()

	err := http.ListenAndServe(":8089", nil)

	if err != nil {
		fmt.Println("Failed to start server: %s", err.Error())
	}
}

