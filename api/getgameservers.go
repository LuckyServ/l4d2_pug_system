package api

import (
	"github.com/gin-gonic/gin"
	"../settings"
	"../utils"
	"strings"
)


func HttpReqGetGameServers(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	mapResponse["success"] = true;
	mapResponse["gameservers"] = settings.GameServers;

	var arIPs []string;
	for _, sIPPORT := range settings.GameServers {
		sIP := strings.Split(sIPPORT, ":")[0];
		if (utils.GetStringIdxInArray(sIP, arIPs) == -1) {
			arIPs = append(arIPs, sIP);
		}
	}

	mapResponse["servers"] = arIPs;
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
