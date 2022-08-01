package api

import (
	"github.com/gin-gonic/gin"
)


func HttpReqPingsReceiver(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	mapResponse["success"] = true;
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
