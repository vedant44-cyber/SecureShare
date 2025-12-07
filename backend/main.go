package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	// AWS SDK v2 (used for MinIO because it speaks S3)
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/credentials" 
)

func main() {
	ctx := context.Background()

	// ---------------- Redis Test ----------------
	fmt.Println("🔍 Checking Redis connection...")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
	})


	err := rdb.Set(ctx, "test-key", "hello-redis", 0).Err()
	if err != nil {
		log.Fatal("❌ Redis connection failed:", err)
	}

	val, _ := rdb.Get(ctx, "test-key").Result()
	fmt.Println("✅ Redis working. Value:", val)

	// ---------------- MinIO Test (S3 compatible) ----------------
	fmt.Println("\n🔍 Checking MinIO connection...")

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               "http://localhost:9000", // local S3 MinIO endpoint (docker service name)
					SigningRegion:     "us-east-1",
					HostnameImmutable: true,
				}, nil
			}),
		),
	)
	if err != nil {
		log.Fatal("❌ Failed to load config:", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		// use credentials package, not aws.*
		o.Credentials = credentials.NewStaticCredentialsProvider(
			"minioadmin",
			"minioadmin123",
			"",
		)
		o.UsePathStyle = true // required for MinIO
	})

	bucket := "secure-share"
	key := "healthcheck.txt"
	content := []byte("MinIO Alive")

	// upload test
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(content),
	})
	if err != nil {
		log.Fatal("❌ MinIO upload failed:", err)
	}

	// download test
	obj, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Fatal("❌ MinIO download failed:", err)
	}

	fmt.Println("✅ MinIO working. File size:", obj.ContentLength)

	/* -----------------------------------------------------------
	   AWS cloud part 
	-------------------------------------------------------------

	// fmt.Println("\n🔍 Checking AWS S3 Production connection...")
	// awsProdCfg, err := config.LoadDefaultConfig(ctx)
	// if err != nil {
	// 	 log.Fatal("❌ AWS config failed:", err)
	// }
	// s3Prod := s3.NewFromConfig(awsProdCfg)
	// _, err = s3Prod.HeadBucket(ctx, &s3.HeadBucketInput{
	// 	 Bucket: aws.String("your-prod-bucket"),
	// })
	// if err != nil {
	// 	 fmt.Println("⚠ AWS not connected yet (expected for dev)")
	// } else {
	// 	 fmt.Println("✅ AWS S3 Production Connected")
	// }

	------------------------------------------------------------- */

	fmt.Println("\n🎉 Redis + MinIO CONNECTED SUCCESSFULLY ")
}
