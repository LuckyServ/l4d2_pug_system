package api

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"io/ioutil"
	"../settings"
	"github.com/gin-contrib/gzip"
)


func GinInit() {
	gin.SetMode(gin.ReleaseMode); //disable debug logs
	gin.DefaultWriter = ioutil.Discard; //disable output
	r := gin.Default();
	r.Use(gzip.Gzip(gzip.DefaultCompression));
	r.MaxMultipartMemory = 1 << 20;

	r.GET("/status", HttpReqStatus);
	r.POST("/shutdown", HttpReqShutdown);
	r.GET("/getme", HttpReqGetMe);
	r.GET("/getonlineplayers", HttpReqGetOnlinePlayers);
	r.GET("/validateprofile", HttpReqValidateProf);
	r.GET("/acceptrules", HttpReqAcceptRules);

	r.GET("/openidcallback", HttpReqOpenID);

	r.GET("/createlobby", HttpReqCreateLobby);
	r.GET("/joinlobby", HttpReqJoinLobby);
	r.GET("/leavelobby", HttpReqLeaveLobby);
	r.GET("/getlobbies", HttpReqGetLobbies);
	r.GET("/joinanylobby", HttpReqJoinAnyLobby);

	r.GET("/myip", HttpReqMyIP);
	
	fmt.Printf("Starting web server\n");
	go r.Run(":"+settings.ListenPort); //Listen on port
}

func HttpReqMyIP(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*");
	c.String(200, c.ClientIP());
}
