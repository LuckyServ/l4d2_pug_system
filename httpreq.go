package main

import (
	"github.com/gin-gonic/gin"
)


func HttpReqPing(c *gin.Context) {
	c.JSON(200, gin.H{"success": true});
}
