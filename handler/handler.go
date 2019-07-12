package handler

import (
	"FILESTORE-SERVER/meta"
	"FILESTORE-SERVER/utils"
	"encoding/json"
	"FILESTORE-SERVER/db"
	"FILESTORE-SERVER/common"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
		r.ParseForm()
		username := r.Form.Get("username")

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

		//更新用户文件
		suc := db.Insert(db.UserFile{
			UserName:username,
			FileSha1:fileMeta.FileSha1,
			FileName:fileMeta.FileName,
			FileSize:fileMeta.FileSize,
			UploadedAt:fileMeta.UpdatedAt,
			UpdatedAt:fileMeta.UpdatedAt,
		})

		fmt.Println("User file insert ",suc)

		fileMetaJson , err := json.Marshal(fileMeta)

		if err != nil {
			fmt.Printf("Failed to json marsha1 err:%s", err.Error())
		}
		io.WriteString(w,fmt.Sprintf("{\"code\":200, \"msg\":success,\"file\":%s}", string(fileMetaJson)))
	}

}

func FastUploadHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	fileSha1 := r.Form.Get("file_sha1")
	fileName := r.Form.Get("file_name")
	fileSize,_:= strconv.ParseInt(r.Form.Get("file_size"),10,64)
	userName := r.Form.Get("username")


	file,err := meta.GetFileDb(fileSha1)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if file == nil {
		w.Write(common.NewResp(-1, "秒传失败,请访问普通接口上传！", nil).JsonBytes())
		return
	}

	//写入文件
	suc := db.Insert(db.UserFile{
		UserName:userName,
		FileName:fileName,
		FileSize:fileSize,
		FileSha1:fileSha1,
	})

	if suc {
		w.Write(common.NewResp(0, "SUCCESS", nil).JsonBytes())
		return
	}else{
		w.Write(common.NewResp(-2, "上传失败,请稍后再试!", nil).JsonBytes())
		return
	}

}


func QueryFile(w http.ResponseWriter, r *http.Request){
	//fileSha1 := r.URL.Query()
	////fileMeta := meta.GetFile(fileSha1["filesha1"][0])
	//file,err := meta.GetFileDb(fileSha1["filesha1"][0])
	//fmt.Println("sha1" + fileSha1["filesha1"][0])
	//if err != nil {
	//	fmt.Printf("Failed to get file meta:%s", err.Error())
	//}
	//
	//fileMeta := meta.FileMeta{
	//	FileName: file.FileName.String,
	//	FileSha1: file.FileSha1.String,
	//	FileSize: file.FileSize.Int64,
	//	FilePath: file.FilePath.String,
	//}
	//
	//fileMetaJson, err := json.Marshal(fileMeta)
	r.ParseForm()

	username := r.Form.Get("username")
	limitStr := r.Form.Get("limit")

	intLimit,_ := strconv.Atoi(limitStr)

	fmt.Println("username:",username)
	fmt.Println("limit:", intLimit)
	ufiles,err := db.QueryUserFile(username, intLimit)

	if err != nil {
		fmt.Println(err.Error())

		resp := common.NewResp(
				int(common.StatusQueryError),
				err.Error(),
				nil,
			)
		w.Write(resp.JsonBytes())
		return
	}

	resp := common.NewResp(
			0,
			"SUCCESS",
			ufiles,
		)
	w.Write(resp.JsonBytes())
}

func SucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Success")
}


