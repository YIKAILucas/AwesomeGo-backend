package oss

import (
	"fmt"
	"github.com/apex/log"
	"github.com/baidubce/bce-sdk-go/services/bos"
)

func main() {
	// 用户的Access Key ID和Secret Access Key
	AK, SK := "ad7910f9ed614f9788d5092ea8c719b0", "4db0368e0b7d4fda800f9ab115638afc"
	bucketName := "tenghui"
	// 用户指定的Endpoint
	ENDPOINT := "tenghui"

	// 初始化一个BosClient
	bosClient, err := bos.NewClient(AK, SK, ENDPOINT)

	err = bosClient.PutBucketAclFromCanned(bucketName, "public-read-write")
	if err != nil {
		log.Error("设置权限失败")
	}

	exists, err := bosClient.DoesBucketExist(bucketName)
	if err == nil && exists {
		fmt.Println("Bucket exists")
	} else {
		fmt.Println("Bucket not exists")
	}

	// 从数据流上传
	//bodyStream, err := bce.NewBodyFromFile(fileName)
	//etag, err := bosClient.PutObject(bucketName, objectName, bodyStream, nil)

	// 从本地文件上传
	fileName := "./baiduyun.go"

	objectName := ""

	etag, err := bosClient.PutObjectFromFile(bucketName, objectName, fileName, nil)

}
