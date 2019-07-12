package handler

import (
	"filestore-server/cache/redis"
	"filestore-server/common"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
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
	fileHash := r.Form.Get("file_hash")
	username := r.Form.Get("username")
	fileSize, _:= strconv.ParseInt(r.Form.Get("file_size"), 10, 64)


	multiInfo := multipartFileInfo{
		FileHash:fileHash,
		FileSize:fileSize,
		UploadId:username + fmt.Sprintf("%x", time.Now().Unix()),
		ChunkCount: int(math.Ceil(float64(fileSize/(5*1024*1024)))),
		ChunkSize:int64(5*1024*1024),
	}

	//3 生成缓存
	redisConn := redis.RedisPool().Get()
	defer redisConn.Close()

	cachePrefix := "MP_"+ multiInfo.UploadId
	redisConn.Do("HSET",cachePrefix,"fileHash",multiInfo.FileHash)
	redisConn.Do("HSET",cachePrefix,"fileSize",multiInfo.FileSize)
	redisConn.Do("HSET",cachePrefix,"chunkSize",multiInfo.ChunkSize)
	redisConn.Do("HSET",cachePrefix,"chunkCount",multiInfo.ChunkCount)

	//4 返回用户
	w.Write(common.NewResp(0,"OK", multiInfo).JsonBytes())
}