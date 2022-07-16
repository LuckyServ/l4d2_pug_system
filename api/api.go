package api

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"io/ioutil"
	"../settings"
)


func GinInit() {
	gin.SetMode(gin.ReleaseMode); //disable debug logs
	gin.DefaultWriter = ioutil.Discard; //disable output
	r := gin.Default();
	r.MaxMultipartMemory = 1 << 20;

	r.GET("/status", HttpReqStatus);
	r.POST("/shutdown", HttpReqShutdown);
	r.GET("/getme", HttpReqGetMe);
	r.GET("/getonlineplayers", HttpReqGetOnlinePlayers);
	r.GET("/openidcallback", HttpReqOpenID);
	r.GET("/validateprofile", HttpReqValidateProf);
	r.GET("/acceptrules", HttpReqAcceptRules);

	r.GET("/createlobby", HttpReqCreateLobby);
	
	fmt.Printf("Starting web server\n");
	go r.Run(":"+settings.ListenPort); //Listen on port
}
