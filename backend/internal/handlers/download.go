package handlers

import (
	"log"
	"secure-share/internal/database"
	"secure-share/internal/helper"
	"secure-share/internal/storage"

	"github.com/gin-gonic/gin"
)

func (depends *HandlerDependencies) HandleDownload(c *gin.Context) {
	// request context
	ctx := c.Request.Context()

	// dependencies
	rdb := depends.RedisClient
	s3 := depends.S3Client
	cfg := depends.Config

	// id sanitization
	id := c.Param("id")
	fileID, err := helper.SanitizeFileID(id)
	if err != nil {
		log.Printf("Invalid file ID: %v", err)
		c.JSON(400, gin.H{"error": "Please provide valid id"})
		return
	}
	// retrieve metadata from Redis
	metaKey := "meta:" + fileID
	metaJSON, err := database.GetMetaJSON(ctx, rdb, metaKey)
	if err != nil {
		if err == database.ErrMetaNotFound {
			c.JSON(404, gin.H{"error": "File not found"})
			return
		}
		log.Printf("Error retrieving metadata: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	//meta is a struct will all the file metadata
	// S3Key,Size ,UploadedAt ,TTLHours,DownloadLimit ,	Filename
	meta, err := helper.ParseFileMetaJSON(metaJSON)
	if err != nil {
		// corrupted metadata treat as expired
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	// retrieve limit from Redis
	limitKey := "limit:" + fileID

	// download limit check
	//case 1 limit key not exist  unlimited downloads
	//case 2 limit key exist  check and decrement

	limitKeyExist := database.LimitExist(ctx, rdb, limitKey)
	if limitKeyExist {
		// limited download logic
		newLimit, err := database.DecrementDownloadLimit(ctx, rdb, limitKey)
		if err != nil {
			log.Printf("Error decrementing download limit for key=%s: %v", limitKey, err)
			// redis error
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		if newLimit < 0 {
			// should not happen due to atomic decrement check
			c.JSON(403, gin.H{"error": "File limit reached"})
			return
		}

		obj, err := storage.GetFileFromS3(ctx, s3, cfg.S3Bucket, meta.S3Key)
		if err != nil {
			log.Printf("Error retrieving file from S3 for key=%s: %v", meta.S3Key, err)
			// s3 error
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		defer obj.Close()
		// stream file to client
		err = helper.StreamFileToClient(c.Writer, obj, meta.Filename, meta.Size)
		if err != nil {
			log.Printf("streaming failed for key=%s: %v", meta.S3Key, err)
		}
		// cleanup after streaming if limit reached
		if newLimit == 0 {
			database.CleanupExpiredKeys(ctx, rdb, metaKey, limitKey)
			storage.CleanupExpiredFiles(ctx, s3, meta.S3Key, cfg.S3Bucket)
		}
		return

	} else {
		// unlimited download logic
		obj, err := storage.GetFileFromS3(ctx, s3, cfg.S3Bucket, meta.S3Key)
		if err != nil {
			log.Printf("Error retrieving file from S3 for key=%s: %v", meta.S3Key, err)
			// s3 error
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		defer obj.Close()
		// stream file to client
		err = helper.StreamFileToClient(c.Writer, obj, meta.Filename, meta.Size)
		if err != nil {
			log.Printf("streaming failed for key=%s: %v", meta.S3Key, err)
		}
		return
	}

}
