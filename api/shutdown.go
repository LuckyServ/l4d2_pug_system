package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
)

var ChShutdown chan bool = make(chan bool);


func HttpReqShutdown(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	if (!auth.Backend(c.PostForm("backend_auth"))) {
		mapResponse["success"] = false;
		c.JSON(200, mapResponse);
		return;
	}

	mapResponse["success"] = true;

	c.Header("Access-Control-Allow-Origin", "*");
	c.JSON(200, mapResponse);
	go PerformShutDown();
}

func PerformShutDown() {
	ChShutdown <- true;
}