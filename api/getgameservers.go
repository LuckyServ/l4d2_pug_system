package api

import (
	"github.com/gin-gonic/gin"
	"../settings"
	"../games"
)


func HttpReqGetGameServers(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	mapResponse["success"] = true;
	games.MuGames.RLock();
	mapResponse["servers"] = settings.GameServers;
	games.MuGames.RUnlock();
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
