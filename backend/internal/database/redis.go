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
