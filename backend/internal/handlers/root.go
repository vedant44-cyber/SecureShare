package handlers

import "github.com/gin-gonic/gin"

func (depends *HandlerDependencies) HandleRoot(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome to the Secure Share API"})
}
