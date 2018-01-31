package main

import (
	"log"
	"strings"

	"github.com/panenming/go-im/libs/contenttype"
	"github.com/panenming/go-im/libs/file"

	"github.com/minio/minio-go"
)

func main() {
	endpoint := "10.39.35.38:9000"
	accessKey := "PYKBI9S2IH1JV00GH5KS"
	secretKey := "ReTLfTEDGI3YWRwD9kbKagUxYFTZTPGkSXpmgQk8"

	useSSL := false

	minioClient, err := minio.New(endpoint, accessKey, secretKey, useSSL)

	if err != nil {
		log.Fatalln(err)
	}

	//log.Printf("%#v\n", minioClient)

	// 创建桶
	bucketName := "media"
	// 所在区，一般是定死的
	location := "z0"

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		exists, err := minioClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	}

	log.Printf("Successfully created %s\n", bucketName)

	// Upload file objectName 唯一
	objectName := "panenming.mp3"
	filePath := "D:/panenming.mp3"
	ext := file.Ext(filePath)
	ext = strings.Trim(ext, ".")
	contentType := contenttype.GetContentTypeByExtension(ext)

	log.Println("contenttype=", contentType)

	n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)

}
