package main

import (
	"FILESTORE-SERVER/handler"
	"fmt"
	"net/http"
)

func main(){
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/success", handler.SucHandler)
	err := http.ListenAndServe(":8089", nil)

	if err != nil {
		fmt.Println("Failed to start server: %s", err.Error())
	}
}

