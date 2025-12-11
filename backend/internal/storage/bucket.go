package storage

import (
	"context"
	"fmt"
	"secure-share/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewS3Client(cfg *config.Config, ctx context.Context) (*minio.Client, error) {

	client, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Secure: cfg.S3UseSSL == "true",
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("S3 is working")
	_, err = client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
