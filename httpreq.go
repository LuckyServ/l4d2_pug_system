package main

import (
	"github.com/gin-gonic/gin"
	"./settings"
	"./auth"
	"io/ioutil"
	"fmt"
	"time"
	"./players"
)


func ginInit() {
	gin.SetMode(gin.ReleaseMode); //disable debug logs
	gin.DefaultWriter = ioutil.Discard; //disable output
	r := gin.Default();
	r.MaxMultipartMemory = 1 << 20;

	r.POST("/status", HttpReqStatus);
	r.POST("/shutdown", HttpReqShutdown);
	r.POST("/addauth", HttpReqAddAuth);
	r.POST("/removeauth", HttpReqRemoveAuth);
	
	fmt.Printf("Starting web server\n");
	go r.Run(":"+settings.ListenPort); //Listen on port
}



func HttpReqStatus(c *gin.Context) {
	mapResponse := map[string]interface{}{
		"success": true,
		"shutdown": bStateShutdown,
		"time": time.Now().UnixMilli(),
		"players_updated": players.I64LastPlayerlistUpdate,
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
	}

	c.JSON(200, mapResponse);
	go PerformShutDown();
}

func HttpReqAddAuth(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	if (bStateShutdown) {
		mapResponse["success"] = false;
		mapResponse["error"] = "Server shutting down";
		c.JSON(200, mapResponse);
		return;
	}
	if (!auth.Backend(c.PostForm("backend_auth"))) {
		mapResponse["success"] = false;
		mapResponse["error"] = "Bad auth key";
		c.JSON(200, mapResponse);
		return;
	}


	sSteamID64 := c.PostForm("steamid64");
	sNicknameBase64 := c.PostForm("nickname_base64");
	if (sSteamID64 == "" || sNicknameBase64 == "") {
		mapResponse["success"] = false;
		mapResponse["error"] = "Not all required parameters present";
		c.JSON(200, mapResponse);
		return;
	}

	sSessionID := auth.AddPlayerAuth(sSteamID64, sNicknameBase64);
	mapResponse["success"] = true;
	mapResponse["session_id"] = sSessionID;
	
	c.JSON(200, mapResponse);
}

func HttpReqRemoveAuth(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	if (bStateShutdown) {
		mapResponse["success"] = false;
		mapResponse["error"] = "Server shutting down";
		c.JSON(200, mapResponse);
		return;
	}
	if (!auth.Backend(c.PostForm("backend_auth"))) {
		mapResponse["success"] = false;
		mapResponse["error"] = "Bad auth key";
		c.JSON(200, mapResponse);
		return;
	}


	sSessID := c.PostForm("session_id");
	if (sSessID == "") {
		mapResponse["success"] = false;
		mapResponse["error"] = "Not all required parameters present";
		c.JSON(200, mapResponse);
		return;
	}

	bSuccess, sError := auth.RemovePlayerAuth(sSessID);
	mapResponse["success"] = bSuccess;
	if (!bSuccess) {
		mapResponse["error"] = sError;
	}
	
	c.JSON(200, mapResponse);
}
