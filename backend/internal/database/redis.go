package database

import (
	"context"
	"errors"
	"fmt"
	"secure-share/internal/config"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var ErrMetaNotFound = errors.New("metadata not found")
var ErrDownloadLimitNotFound = errors.New("download not found")

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

func DeleteFileMetadata(ctx context.Context, rdb *redis.Client, metaKey, limitKey string) error {
	_, err := rdb.Del(ctx, metaKey, limitKey).Result()
	return err
}

func DecrementDownloadLimit(ctx context.Context, rdb *redis.Client, limitKey string) (int, error) {
	val, err := rdb.Decr(ctx, limitKey).Result()
	if err != nil {
		return 0, err // real Redis error
	}
	return int(val), nil
}

func UpdateDownloadLimit(ctx context.Context, rdb *redis.Client, limitKey string, newLimit int) error {
	_, err := rdb.Set(ctx, limitKey, newLimit, 0).Result()
	return err
}

func GetDownloadLimit(ctx context.Context, rdb *redis.Client, limitKey string) (int, error) {
	val, err := rdb.Get(ctx, limitKey).Result()
	if err == redis.Nil {
		return 0, ErrDownloadLimitNotFound // no limit set
	}
	if err != nil {
		return 0, err //  Redis error
	}
	limit, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return limit, nil
}

func CleanupExpiredKeys(ctx context.Context, rdb *redis.Client, metaKey, limitKey string) {
	if err := rdb.Del(ctx, limitKey).Err(); err != nil {
		fmt.Printf("cleanup: failed to delete Redis key %s : %v\n", limitKey, err)
	}
	if err := rdb.Del(ctx, metaKey).Err(); err != nil {
		fmt.Printf("cleanup: failed to delete Redis key %s: %v\n", metaKey, err)
	}
}

func LimitExist(ctx context.Context, rdb *redis.Client, limitKey string) bool {
	val, err := rdb.Exists(ctx, limitKey).Result()
	if err != nil {
		return false
	}
	return val == 1
}
