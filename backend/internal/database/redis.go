package database

import (
	"context"
	"fmt"
	"secure-share/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config, ctx context.Context) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass, // no password
		DB:       0,             // use default DB
		Protocol: 2,
	})
	fmt.Println("Redis is working")
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}
