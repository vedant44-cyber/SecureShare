package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// Implementation for handling file upload will go here
// input parametes -TTL(integer) in sec,Download_limit (integer),file(binaryFile),filename(string),
func (depends *HandlerDependencies) HandleUpload(c *gin.Context) {
	ctx := c.Request.Context()

	//  Parse TTL (optional)
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
		if err != nil || n < 0 {
			c.JSON(400, gin.H{"error": "invalid download_limit"})
			return
		}
		downloadLimit = n
	}

	filename := c.PostForm("filename")
	if filename == "" {
		c.JSON(400, gin.H{"error": "filename is required"})
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
		c.JSON(500, gin.H{"error": "s3 upload failed", "detail": err.Error()})
		return
	}

	//  Write metadata to Redis
	metaKey := "meta:" + fileID
	limitKey := "limit:" + fileID

	uploadedAt := time.Now().Unix()

	// TTL duration
	var expire time.Duration
	if ttl > 0 {
		expire = time.Duration(ttl) * time.Second
	} else {
		expire = 0
	}

	// Store meta data as JSON string
	meta := fmt.Sprintf(
		`{"s3_key":"%s","size":%d,"uploaded_at":%d,"ttl_seconds":%d,"download_limit":%d,"filename":"%s"}`,
		s3Key,
		uploadInfo.Size,
		uploadedAt,
		ttl,
		downloadLimit,
		filename,
	)

	err = depends.RedisClient.Set(ctx, metaKey, meta, expire).Err()
	if err != nil {
		// rollback S3 object
		depends.S3Client.RemoveObject(ctx, depends.Config.S3Bucket, s3Key, minio.RemoveObjectOptions{})
		c.JSON(500, gin.H{"error": "metadata write failed"})
		return
	}

	// Store download limit (if > 0)
	if downloadLimit > 0 {
		err = depends.RedisClient.Set(ctx, limitKey, downloadLimit, expire).Err()
		if err != nil {
			// rollback both
			depends.RedisClient.Del(ctx, metaKey)
			depends.S3Client.RemoveObject(ctx, depends.Config.S3Bucket, s3Key, minio.RemoveObjectOptions{})
			c.JSON(500, gin.H{"error": "limit write failed"})
			return
		}
	}

	// Return response
	c.JSON(200, gin.H{
		"file_id":        fileID,
		"size":           uploadInfo.Size,
		"ttl_seconds":    ttl,
		"download_limit": downloadLimit,
		"filename":       filename,
	})
}
