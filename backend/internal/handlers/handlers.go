package handlers

import (
	"secure-share/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type HandlerDependencies struct {
	RedisClient *redis.Client
	S3Client    *minio.Client
	Config      *config.Config
}
