package handlers

import "github.com/gin-gonic/gin"

func (depends *HandlerDependencies) HandleDownload(c *gin.Context) {
	c.JSON(200, gin.H{"message": "ok"})
}
