package database

import (
	"context"
	"errors"
	"fmt"
	"secure-share/internal/config"

	"github.com/redis/go-redis/v9"
)

var ErrMetaNotFound = errors.New("metadata not found")

func NewRedisClient(cfg *config.Config, ctx context.Context) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       0,
		Protocol: 2,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	fmt.Println("Connected to Redis successfully")
	return client, nil
}
func GetMetaJSON(ctx context.Context, rdb *redis.Client, metaKey string) (string, error) {
	val, err := rdb.Get(ctx, metaKey).Result()
	if err == redis.Nil {
		return "", ErrMetaNotFound //metadata  expired or deleted
	}

	if err != nil {
		return "", err // real Redis error
	}
	return val, nil
}
