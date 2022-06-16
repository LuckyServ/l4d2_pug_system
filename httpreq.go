package main

import (
	"github.com/gin-gonic/gin"
	"./settings"
	"./auth"
	"io/ioutil"
	"fmt"
)


func ginInit() {
	gin.SetMode(gin.ReleaseMode); //disable debug logs
	gin.DefaultWriter = ioutil.Discard; //disable output
	r := gin.Default();
	r.MaxMultipartMemory = 1 << 20;

	r.POST("/status", HttpReqStatus);
	r.POST("/shutdown", HttpReqShutdown);
	
	fmt.Printf("Starting web server\n");
	go r.Run(":"+settings.ListenPort); //Listen on port
}



func HttpReqStatus(c *gin.Context) {
	mapResponse := map[string]interface{}{
		"success": true,
		"shutdown": bStateShutdown,
	};
	c.JSON(200, mapResponse);
}

func HttpReqShutdown(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	if (!auth.Backend(c.PostForm("backend_auth"))) {
		mapResponse["success"] = false;
		mapResponse["error"] = "Bad auth key";
		c.JSON(200, mapResponse);
		return;
	}

	bSetShutdown, sError := SetShutDown();
	if (!bSetShutdown) {
		mapResponse["success"] = false;
		mapResponse["error"] = sError;
	} else {
		mapResponse["success"] = true;
		go PerformShutDown();
	}

	c.JSON(200, mapResponse);
}
