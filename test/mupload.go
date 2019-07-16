package main

import (
	"bufio"
	"bytes"
	"fmt"
	jsonit "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func multiupload(filename ,targeUrl string,chunkSize int) error{
	fp,err := os.Open(filename)

	if err != nil {
		fmt.Println("Open file err", err.Error())
		return err
	}

	indexChan := make(chan int)

	index := 0
	freader := bufio.NewReader(fp)

	buf := make([]byte, chunkSize)
	for {
		n,err := freader.Read(buf)

		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			break
		}

		index ++

		bufCopied := make([]byte, 5*1024*1024)
		copy(bufCopied, buf)
		//fp ,_ := os.Create(strconv.Itoa(index))
		//defer fp.Close()

		//fp.Write(bufCopied[:n])
		go func(fbuf []byte,curIndex int){

			resp, err := http.Post(
				targeUrl + "&index=" + strconv.Itoa(curIndex),
				"multipart/form-data",
				bytes.NewReader(fbuf),
			)

			if err != nil {
				fmt.Println(err.Error())
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			fmt.Println("index:%d,body:%s", curIndex,string(body))
			indexChan <- curIndex
		}(bufCopied[:n],index)
	}

	for i:= 0 ; i < index; i++ {
		select{
			case res := <- indexChan:
				fmt.Println(res)
		}
	}
	//fp1 ,_ := os.Create("vs.exe")
	//defer fp1.Close()


	//for i:=1; i <= 11; i++{
	//	fchunk,_:= os.Open(strconv.Itoa(i))
	//	offset := (i-1)*5*1024*1024
	//
	//	n,_:= fchunk.Read(buf)
	//
	//	fp1.WriteAt(buf[:n], int64(offset))
	//}
	return nil
}
func main(){
	//参数
	username := "rains"
	token := "648f21ad845778e36a7a2a4808b8b5705d2c29fb"
	fileHash := "705b1098c525406e48d3fe3c02373ff380cc97d5"
	fileSize := "52918224"
	baseUrl := "http://localhost:8089"
	//上传
	urlInit := baseUrl + "/file/initupload"
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
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil{
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	fmt.Println(string(body))
	//解析参数
	uploadId := jsonit.Get(body,"Data").Get("UploadId").ToString()
	chunkSize := jsonit.Get(body,"Data").Get("ChunkSize").ToInt()

	fmt.Printf("uploadId:%s,chunkSize:%d", uploadId, chunkSize)

	//多协程上传文件
	filename := "d:\\VSCodeUserSetup-x64-1.36.0.exe"

	urlMupload := baseUrl + "/file/mupload?uploadid=" + uploadId + "&username="+username + "&token=" + token
	err = multiupload(filename,urlMupload,chunkSize)

	if err != nil{
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	//上传查询
	urlComplete := baseUrl + "/file/completeupload"

	resp, err = http.PostForm(
		urlComplete,
		url.Values{
			"username":{username},
			"token":{token},
			"uploadid":{uploadId},
			"filehash":{fileHash},
			"filename":{"vs.exe"},
		},
	)

	if err != nil {
		fmt.Println(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}