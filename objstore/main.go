package main

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	s3Client, err := minio.New("127.0.0.1:9000", &minio.Options{
		Creds: credentials.NewStaticV4("Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", ""),
	})
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := s3Client.FPutObject(context.Background(), "tmp", "go.sum", "./go.sum", minio.PutObjectOptions{
		ContentType: "text",
	}); err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully uploaded")
}
