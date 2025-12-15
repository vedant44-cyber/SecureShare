package handlers

import (
	"log"
	"strconv"
	"time"

	"secure-share/internal/helper"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// Implementation for handling file upload will go here
// input parameters -TTL(string) in hour,Download_limit (string),file(file),filename(string),
func (depends *HandlerDependencies) HandleUpload(c *gin.Context) {
	ctx := c.Request.Context()

	//  Parse TTL parameter
	ttlStr := c.PostForm("ttl")
	ttl := 0
	if ttlStr != "" {
		n, err := strconv.Atoi(ttlStr)
		if err != nil || n < 0 {
			c.JSON(400, gin.H{"error": "invalid ttl"})
			return
		}
		ttl = n
	}

	// Parse download_limit
	limitStr := c.PostForm("download_limit")
	downloadLimit := 0
	if limitStr != "" {
		n, err := strconv.Atoi(limitStr)
		// download limit 1 hour to 30 days
		if err != nil || n < 0 || n > 720 {
			c.JSON(400, gin.H{"error": "invalid download_limit"})
			return
		}
		downloadLimit = n
	}

	rawfilename := c.PostForm("filename")
	if rawfilename == "" {
		c.JSON(400, gin.H{"error": "filename is required"})
		return
	}
	filename := helper.SanitizeFilename(rawfilename)
	if filename == "file" {
		c.JSON(400, gin.H{"error": "invalid filename"})
		return
	}

	//  Get encrypted blob from frontend
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "file is required"})
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer src.Close() //close the file after function ends

	//  Create fileID + s3_key
	fileID := uuid.New().String()
	s3Key := "files/" + fileID

	// Upload encrypted blob to S3
	uploadInfo, err := depends.S3Client.PutObject(ctx,
		depends.Config.S3Bucket,
		s3Key,
		src,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		},
	)

	if err != nil {
		log.Printf("s3 upload failed for key=%s: %v", s3Key, err)
		c.JSON(503, gin.H{"error": "Service temporarily unavailable. Please retry."})
		return
	}

	//  Write metadata to Redis
	metaKey := "meta:" + fileID
	limitKey := "limit:" + fileID

	uploadedAt := time.Now().Unix()

	// TTL duration in hours
	var expire time.Duration
	if ttl > 0 {
		expire = time.Duration(ttl) * time.Hour
	} else {
		expire = 0
	}

	// Store meta data as JSON string
	meta, err := helper.BuildFileMetaJSON(s3Key, uploadInfo.Size, uploadedAt, ttl, downloadLimit, filename)
	if err != nil {
		// rollback S3 object
		if removeErr := depends.S3Client.RemoveObject(ctx, depends.Config.S3Bucket, s3Key, minio.RemoveObjectOptions{}); removeErr != nil {
			log.Printf("rollback failed: unable to remove s3 object key=%s: %v", s3Key, removeErr)
		}
		c.JSON(500, gin.H{"error": "Upload failed. Please try again later."})
		log.Printf("failed to build file metadata for metaKey=%s: %v", metaKey, err)
		return
	}
	// write to redis
	err = depends.RedisClient.Set(ctx, metaKey, meta, expire).Err()
	if err != nil {
		// rollback S3 object
		if removeErr := depends.S3Client.RemoveObject(ctx, depends.Config.S3Bucket, s3Key, minio.RemoveObjectOptions{}); removeErr != nil {
			log.Printf("rollback failed: unable to remove s3 object key=%s: %v", s3Key, removeErr)
		}
		c.JSON(503, gin.H{"error": "Service temporarily unavailable. Please retry."})
		log.Printf("redis SET failed for metaKey=%s: %v", metaKey, err)
		return
	}
	// download_limit =0 means unlimited downloads
	// Store download limit (if > 0)
	if downloadLimit > 0 {
		err = depends.RedisClient.Set(ctx, limitKey, downloadLimit, expire).Err()
		if err != nil {
			// rollback both
			if delErr := depends.RedisClient.Del(ctx, metaKey).Err(); delErr != nil {
				log.Printf("rollback failed: unable to delete redis key=%s: %v", metaKey, delErr)
			}
			if removeErr := depends.S3Client.RemoveObject(ctx, depends.Config.S3Bucket, s3Key, minio.RemoveObjectOptions{}); removeErr != nil {
				log.Printf("rollback failed: unable to remove s3 object key=%s: %v", s3Key, removeErr)
			}
			c.JSON(500, gin.H{"error": "Upload failed. Please try again later."})
			return
		}
	}

	// Return response
	response := helper.BuildUploadResponse(fileID, uploadInfo.Size, ttl, downloadLimit, filename)
	c.JSON(201, response)
	// resource created successfully
}
