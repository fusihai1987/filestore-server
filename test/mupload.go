package main

import (
	jsonit "github.com/json-iterator/go"
	"net/http"
	"net/url"
	"fmt"
	"os"
	"io/ioutil"
)
func main(){
	//参数
	username := "rains"
	token := "d67a352419f73b2eaf6ec3013e1ac5065d2a945b"
	fileHash := "020a5a011dba5dbe4581f5e66b4e00cfcfb50e2b"
	fileSize := "34449"
	//上传
	urlInit := "http://localhost:8089/file/mupload"
	resp, err := http.PostForm(
		urlInit,
		url.Values{
			"username":{username},
			"token":{token},
			"filehash":{fileHash},
			"filesize":{fileSize},
		},
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer resp.Body.Close()
	body, err :=ioutil.ReadAll(resp.Body)

	if err != nil{
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	fmt.Println(string(body))
	//解析参数
	uploadId := jsonit.Get(body,"Data").Get("UploadId").ToString()
	chunkSize := jsonit.Get(body,"Data").Get("ChunkSize").ToString()

	fmt.Printf("uploadId:%s,chunkSize:%s", uploadId,chunkSize)
}