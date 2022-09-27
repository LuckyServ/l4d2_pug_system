package api

import (
	"github.com/gin-gonic/gin"
	"../settings"
)


func HttpReqGetGameServers(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	mapResponse["success"] = true;
	mapResponse["servers"] = settings.GameServers;
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
