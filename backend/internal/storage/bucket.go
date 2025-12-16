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

	_, err = client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to S3 storage successfully")
	return client, nil
}
func FileExists(ctx context.Context, s3 *minio.Client, S3Bucket string, filename string) bool {

	_, err := s3.StatObject(ctx, S3Bucket, filename, minio.StatObjectOptions{})
	return err == nil

}
func DeleteFile(ctx context.Context, s3 *minio.Client, S3Bucket string, objectName string) error {
	err := s3.RemoveObject(ctx, S3Bucket, objectName, minio.RemoveObjectOptions{})
	return err
}
