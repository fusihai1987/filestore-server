package main

import (
	"filestore-server/store/ceph"
	"fmt"
	"gopkg.in/amz.v1/s3"
	"os"
)

func main(){
	bucket := ceph.GetCephBucket("bucktest1")
	//up_test(bucket)
	download_test(bucket)
}

func up_test(bucket *s3.Bucket){
	//创建bucket
	err := bucket.PutBucket(s3.PublicRead)
	fmt.Printf("create bucktet error:v%\n",err)

	//查询bucket下的keys
	res, _ := bucket.List("","","",100)
	fmt.Printf("object keys:v%\n", res)

	//上传文件测试
	err = bucket.Put("ceph_key.txt",[]byte("just for test!"),"octet-stream", s3.PublicRead)
	fmt.Printf("upload err:%v\n",err)

	//查询bucket下的keys
	res, _ = bucket.List("","","",100)
	fmt.Printf("object keys:v%\n", res)
}

func download_test(bucket *s3.Bucket){

	d, _ := bucket.Get("ceph_key.txt")
	tmpFile, _ := os.Create("test_file.txt")
	tmpFile.Write(d)
	return
}

func test(){
}
