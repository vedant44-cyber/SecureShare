package main

import (
	"context"
	"fmt"
	"secure-share/internal/config"
	"secure-share/internal/database"
	"secure-share/internal/handlers"
	"secure-share/internal/helper"
	"secure-share/internal/router"
	"secure-share/internal/storage"
)

func main() {
	//context
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	helper.ErrorHandler(err)
	//fmt.Println("Starting the application\n", cfg)

	// redisclient
	redisclient, err := database.NewRedisClient(cfg, ctx)
	helper.ErrorHandler(err)
	//fmt.Println("starting the redis\n", redisclient, ctx)

	// s3client
	s3client, err := storage.NewS3Client(cfg, ctx)
	helper.ErrorHandler(err)
	//fmt.Println("starting the s3\n", s3client, ctx)

	dependency := &handlers.HandlerDependencies{
		RedisClient: redisclient,
		S3Client:    s3client,
		Config:      cfg,
	}

	//router
	r, err := router.NewRouter(dependency)
	helper.ErrorHandler(err)
	fmt.Println(r)

	//start server

}
