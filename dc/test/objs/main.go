package main

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	s3Client, err := minio.New("127.0.0.1:9000", &minio.Options{
		Creds: credentials.NewStaticV4("liuz0063", "12345678", ""),
	})
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := s3Client.FPutObject(context.Background(), "tmp", "ok.png", "./test/assets/img.png", minio.PutObjectOptions{}); err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully uploaded")
}
