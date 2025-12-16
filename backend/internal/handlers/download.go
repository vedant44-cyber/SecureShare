package handlers

import (
	"log"
	"secure-share/internal/database"
	"secure-share/internal/helper"

	"github.com/gin-gonic/gin"
)

//TODO
//
//
//
//
//
//
//
//

func (depends *HandlerDependencies) HandleDownload(c *gin.Context) {
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

	// download limit check
	//case 1 limit key not exist  unlimited downloads
	//case 2 limit key exist  check and decrement
	limitKey := "limit:" + fileID
}
