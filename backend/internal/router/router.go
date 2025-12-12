package router

import (
	"secure-share/internal/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter(dependency *handlers.HandlerDependencies) (*gin.Engine, error) {
	// Router initialization logic goes here
	r := gin.New()
	r.GET("/", dependency.HandleRoot)
	r.POST("/upload", dependency.HandleUpload)
	r.GET("/download/:id", dependency.HandleDownload)
	return r, nil
}

