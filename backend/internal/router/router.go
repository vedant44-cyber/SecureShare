package router

import (
	"os"
	"strings"

	"secure-share/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(dependency *handlers.HandlerDependencies) (*gin.Engine, error) {

	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if err := r.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	allowedOrigins := os.Getenv("CORS_ORIGINS")
	origins := []string{}
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Disposition",
			"Content-Length",
			"x-iv",
			"x-filename",
		},
		AllowCredentials: false,
	}))

	r.GET("/", dependency.HandleRoot)
	r.POST("/upload", dependency.HandleUpload)
	r.GET("/download/:id", dependency.HandleDownload)

	return r, nil
}
