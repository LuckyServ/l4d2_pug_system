package main

import (
	"github.com/gin-gonic/gin"
	"./settings"
	"io/ioutil"
	"fmt"
)


func ginInit() {
	gin.SetMode(gin.ReleaseMode); //disable debug logs
	gin.DefaultWriter = ioutil.Discard; //disable output
	r := gin.Default();
	r.MaxMultipartMemory = 1 << 20;

	r.GET("/ping", HttpReqPing);
	
	fmt.Printf("Starting web server\n");
	go r.Run(":"+settings.ListenPort); //Listen on port
}



func HttpReqPing(c *gin.Context) {
	c.JSON(200, gin.H{"success": true});
}
