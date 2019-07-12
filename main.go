package main

import (
	"FILESTORE-SERVER/handler"
	"fmt"
	"net/http"
)

func main(){
	http.Handle("/static/",
		http.StripPrefix("/static/",http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/file/upload", handler.HttpMiddle(handler.UploadHandler))
	http.HandleFunc("/file/success", handler.SucHandler)
	http.HandleFunc("/file/query",handler.QueryFile)
	http.HandleFunc("/user/signup",  handler.SignUpHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.HttpMiddle(handler.QueryUserInfo))
	http.HandleFunc("/file/fastupload", handler.HttpMiddle(handler.FastUploadHandler))

	err := http.ListenAndServe(":8089", nil)

	if err != nil {
		fmt.Println("Failed to start server: %s", err.Error())
	}
}

