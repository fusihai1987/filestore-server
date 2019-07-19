package ceph

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"

	cfg "filestore-server/config"
)

var cephConn *s3.S3

func GetCephConn() *s3.S3 {

	if cephConn != nil {
		return cephConn
	}

	// auth认证
	auth := aws.Auth{
		AccessKey:cfg.CephAccessKey,
		SecretKey:cfg.CephSecretKey,
	}

	curRegion := aws.Region{
		Name:"default",
		EC2Endpoint:cfg.CephEendpoint,
		S3Endpoint:cfg.CephEendpoint,
		S3BucketEndpoint:"",
		S3LocationConstraint:false,
		S3LowercaseBucket:false,
		Sign:aws.SignV2,
	}

	return s3.New(auth, curRegion)
}

//获取指定bucket对象
func GetCephBucket(bucket string) *s3.Bucket {
	conn := GetCephConn()
	return conn.Bucket(bucket)
}

//上传对象
func PutObject(bucket string, path string, data []byte) error {
	return GetCephBucket(bucket).Put(path,data,"octet-stream",s3.PublicRead)
}








