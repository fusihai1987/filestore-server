package handler

import (
	"filestore-server/cache/redis"
	redigo "github.com/garyburd/redigo/redis"
	"filestore-server/common"
	"fmt"
	"math"
	"net/http"
	"path"
	"sort"
	"strconv"
	"time"
	"os"
	"strings"
	"filestore-server/meta"
	"filestore-server/db"
)

type multipartFileInfo struct {
	FileHash string
	FileSize int64
	UploadId string
	ChunkSize int64
	ChunkCount int
}


func InitUploadInfo(w http.ResponseWriter, r *http.Request){
	//1 解析前端参数
	r.ParseForm()
	fileHash := r.Form.Get("filehash")
	username := r.Form.Get("username")
	fileSize, _:= strconv.ParseInt(r.Form.Get("filesize"), 10, 64)



	multiInfo := multipartFileInfo{
		FileHash:fileHash,
		FileSize:fileSize,
		UploadId:username + fmt.Sprintf("%x", time.Now().Unix()),
		ChunkCount: int(math.Ceil(float64(float32(fileSize)/float32((5*1024*1024))))),
		ChunkSize: int64(5*1024*1024),
	}
	fmt.Println("uploadId",multiInfo.UploadId)
	fmt.Println("chunksize",multiInfo.ChunkSize)
	fmt.Println("chunkcount", multiInfo.ChunkCount)
	//3 生成缓存
	redisConn := redis.RedisPool().Get()
	defer redisConn.Close()

	cachePrefix := "MP_"+ multiInfo.UploadId
	redisConn.Do("HSET",cachePrefix,"fileHash",multiInfo.FileHash)
	redisConn.Do("HSET",cachePrefix,"fileSize",multiInfo.FileSize)
	redisConn.Do("HSET",cachePrefix,"chunkSize",multiInfo.ChunkSize)
	redisConn.Do("HSET",cachePrefix,"chunkCount",multiInfo.ChunkCount)

	//4 返回用户
	resp := common.NewResp(0,"OK", multiInfo)
	w.Write(resp.JsonBytes())
}
func MartiUploadHandle(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	//1 解析参数
	chunkIndex := r.Form.Get("index")
	uploadId := r.Form.Get("uploadid")

	//2 创建目录
	fpath := "/d/" + uploadId +"/" + chunkIndex
	err := os.MkdirAll(path.Dir(fpath), 0744)
	if err != nil {
		w.Write(common.NewResp(-1,"mkdir err", nil).JsonBytes())
		return
	}
	fp, err := os.Create(fpath)

	if err!= nil {
		w.Write(common.NewResp(-1,"create err" + err.Error(), nil).JsonBytes())
		return
	}
	defer fp.Close()

	buf := make([]byte,1024*1024)
	offset := int64(0)
	for {
		n, err := r.Body.Read(buf)

		fp.WriteAt(buf[:n], int64(offset))
		offset += int64(n)
		if err != nil {
			break
		}

	}

	//3 缓存index
	redisConn := redis.RedisPool().Get()
	defer redisConn.Close()
	redisConn.Do("HSET","MP_"+uploadId,"CHUNK_INDEX"+chunkIndex,1)

	//4 返回参数
	w.Write(common.NewResp(0,"OK",nil).JsonBytes())
}

func CompeteUploadHandler(w http.ResponseWriter,r *http.Request){
	// 1 解析参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	uploadid := r.Form.Get("uploadid")
	filename := r.Form.Get("filename")
	filesize,_ := strconv.ParseInt(r.Form.Get("filesize"), 10, 64)
	// 2 对比缓存
	redisConn := redis.RedisPool().Get()
	defer redisConn.Close()

	data,err := redigo.Values(redisConn.Do("HGETALL","MP_"+uploadid))

	if err != nil {
		w.Write(common.NewResp(-1,"Context invalid", nil).JsonBytes())
		return
	}

	totalChunkCnt := 0
	chunkCnt := 0
	chunkFiles := make([]int,0)
	for i:=0; i < len(data); i+=2 {
		k := string(data[i].([]byte))
		v := string(data[i + 1].([]byte))
		fmt.Printf("K:%s,V:%s\n", k, v)
		if k == "chunkCount"{
			totalChunkCnt, _= strconv.Atoi(v)
		}else if(strings.HasPrefix(k,"CHUNK_INDEX") && v == "1"){
			chunkCnt += 1
			file , _ := strconv.Atoi(strings.TrimPrefix(k,"CHUNK_INDEX"))
			chunkFiles = append(chunkFiles, file)
		}

	}

	if totalChunkCnt != chunkCnt {
		w.Write(common.NewResp(-2,"Invalid parameter",nil).JsonBytes())
		return
	}
	// 3 合并文件
	path := "/d/" + uploadid
	err = mergeFiles(path, chunkFiles, filename)

	if err != nil {
		w.Write(common.NewResp(-1,fmt.Sprintf("文件合并失败:%s",err.Error()),nil).JsonBytes())
		return
	}
	// 4 写入文件和用户文件
	_ = meta.Insert(meta.FileMeta{FileName:filename,FileSha1:filehash,FileSize:filesize})
	_ = db.Insert(db.UserFile{UserName:username,FileSha1:filehash,FileSize:filesize,FileName:filename,
		UploadedAt:time.Now().Format("2006-1-2 15:04:05"),UpdatedAt:time.Now().Format("2006-1-2 15:04:05")})
	
	// 5
	w.Write(common.NewResp(0,"OK",nil).JsonBytes())
}

func mergeFiles(path string,files []int, targetFileName string) error{

	targetFile,err := os.Create(path + "/" + targetFileName)
	fmt.Println("path:",path + "/" + targetFileName)
	defer targetFile.Close()
	if err != nil {
		return err
	}

	sort.Ints(files)
	buf := make([]byte, 5*1024*1024)
	offset := 0
	for _, file := range files {
		chunFile := path + "/" + strconv.Itoa(file)
		fchunk, err := os.Open(chunFile)
		defer fchunk.Close()

		if err != nil {
			return err
		}
		n,_ := fchunk.Read(buf)
		targetFile.WriteAt(buf[:n], int64(offset))

		offset += n
	}

	return nil
}

