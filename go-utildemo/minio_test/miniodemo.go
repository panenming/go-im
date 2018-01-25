package main

import (
	"log"

	"github.com/minio/minio-go"
)

func main() {
	endpoint := "127.0.0.1:9000"
	accessKey := "OV1CZC6UHXEFCVWEXWAO"
	secretKey := "Om9BnBoeMqTCVvAkuJGpYMyddIsimHIvbZ70Rywn"

	useSSL := false

	minioClient, err := minio.New(endpoint, accessKey, secretKey, useSSL)

	if err != nil {
		log.Fatalln(err)
	}

	//log.Printf("%#v\n", minioClient)

	// 创建桶
	bucketName := "gofile"
	// 所在区，一般是定死的
	location := "us-east-1"

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

	// Upload file
	objectName := "Windows.zip"
	filePath := "D:/Windows.zip"
	contentType := "application/zip"

	n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)

}
