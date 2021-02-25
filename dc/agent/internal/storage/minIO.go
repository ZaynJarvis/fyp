package storage

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIO struct {
	*minio.Client
	root string
}

func newMinIO(addr string) (*MinIO, error) {
	s3Client, err := minio.New(addr, &minio.Options{
		Creds: credentials.NewStaticV4("Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", ""),
	})
	if err != nil {
		log.Fatalln(err)
	}
	tmpDir, err := ioutil.TempDir("", "minIO")
	if err != nil {
		return nil, err
	}
	return &MinIO{
		Client: s3Client,
		root:   tmpDir,
	}, nil
}

func (m *MinIO) Close() {
}

func (m *MinIO) Image(id, contentType string, data []byte) error {
	fn := path.Join(m.root, id)
	if _, err := os.Stat(id); os.IsExist(err) {
		return errors.New("file exists")
	}
	if err := ioutil.WriteFile(fn, data, os.ModePerm); err != nil {
		return errors.New("write file to temporary location failed")
	}
	bucket := time.Now().Format("2006-01-02")
	ctx := context.Background()
	exists, err := m.BucketExists(ctx, bucket)
	if err != nil {
		logrus.Error(err)
	} else if !exists {
		if err := m.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
			Region: "sg",
		}); err != nil {
			logrus.Error("make bucket error: ", err)
		}
	}

	if _, err := m.FPutObject(ctx, bucket, id, fn, minio.PutObjectOptions{
		ContentType: contentType,
	}); err != nil {
		return errors.New("file upload to cloud storage failed")
	}
	logrus.Debug("Successfully uploaded")
	return nil
}
