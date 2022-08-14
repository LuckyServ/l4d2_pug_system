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

	r.GET("/shutdown", HttpReqShutdown);

	r.GET("/status", HttpReqStatus);
	r.GET("/validateprofile", HttpReqValidateProf);
	r.GET("/acceptrules", HttpReqAcceptRules);

	r.GET("/getonlineplayers", HttpReqGetOnlinePlayers);

	r.GET("/createlobby", HttpReqCreateLobby);
	r.GET("/joinlobby", HttpReqJoinLobby);
	r.GET("/leavelobby", HttpReqLeaveLobby);
	r.GET("/getlobbies", HttpReqGetLobbies);
	r.GET("/joinanylobby", HttpReqJoinAnyLobby);
	r.GET("/readyup", HttpReqReadyUp);

	r.GET("/openidcallback", HttpReqOpenID);
	r.GET("/myip", HttpReqMyIP);
	r.GET("/home", HttpReqHome);
	r.POST("/home", HttpReqHome);

	r.GET("/getgame", HttpReqGetGame);

	r.GET("/getgameservers", HttpReqGetGameServers);
	r.GET("/pingsreceiver", HttpReqPingsReceiver);


	r.POST("/gs/getgame", HttpReqGSGetGame);
	r.POST("/gs/fullrup", HttpReqGSFullReadyUp);

	
	fmt.Printf("Starting web server\n");
	go r.Run(":"+settings.ListenPort); //Listen on port
}

func HttpReqMyIP(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.String(200, c.ClientIP());
}

func HttpReqHome(c *gin.Context) {
	sLink := c.Query("link");

	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	if (sLink == "") {
		c.Redirect(303, "https://"+settings.HomeDomain);
	} else {
		c.Redirect(303, sLink);
	}
}
