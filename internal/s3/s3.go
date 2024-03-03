package s3

import (
	"context"
	"io"
	"log"
	"quree/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClientExtended struct {
	*minio.Client
	Ctx        context.Context
	BucketName string
}

type S3Config struct {
	Ctx        context.Context
	Endpoint   string
	AccessKey  string
	SecretKey  string
	UseSSL     bool
	BucketName string
}

var S3Client = Init(&S3Config{
	Ctx:        context.Background(),
	Endpoint:   config.S3_ENDPOINT,
	BucketName: config.S3_BUCKET,
	AccessKey:  config.S3_ACCESS_KEY,
	SecretKey:  config.S3_SECRET_KEY,
	UseSSL:     false,
})

func Init(c *S3Config) *MinioClientExtended {
	client, err := minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKey, c.SecretKey, ""),
		Secure: c.UseSSL,
	})
	if err != nil {
		log.Fatal("Error creating S3Client:", err)
	}

	clientExtended := &MinioClientExtended{
		Client:     client,
		Ctx:        c.Ctx,
		BucketName: c.BucketName,
	}

	return clientExtended

}

func (mce *MinioClientExtended) UploadImage(objectName string, reader io.Reader, objectSize int64) (minio.UploadInfo, error) {
	info, err := mce.PutObject(mce.Ctx, mce.BucketName, objectName, reader, objectSize, minio.PutObjectOptions{})
	if err != nil {
		return info, err
	}

	return info, nil
}
