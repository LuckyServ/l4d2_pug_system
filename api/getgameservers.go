package api

import (
	"github.com/gin-gonic/gin"
	"../settings"
)


func HttpReqGetGameServers(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	mapResponse["success"] = true;
	mapResponse["gameservers"] = settings.GameServers;
	mapResponse["servers"] = settings.HardwareServers;
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
