package api

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"io/ioutil"
	"../settings"
	"github.com/remerge/gzip"
)


func GinInit() {
	gin.SetMode(gin.ReleaseMode); //disable debug logs
	gin.DefaultWriter = ioutil.Discard; //disable output
	r := gin.Default();
	r.Use(gzip.Gzip(gzip.DefaultCompression));
	r.MaxMultipartMemory = 1 << 20;

	r.GET("/shutdown", HttpReqShutdown);
	r.GET("/blocknewgames", HttpReqBlockNewGames);
	r.GET("/refreshservers", HttpRefreshServers);

	r.GET("/ws", HttpReqWebSocket);

	r.GET("/addban", HttpReqAddBan);
	r.GET("/unban", HttpReqUnban);
	r.GET("/deleteban", HttpReqDeleteBan);
	r.GET("/setadmin", HttpReqSetAdmin);
	r.GET("/setmmr", HttpReqSetMmr);
	r.GET("/getknownaccs", HttpReqGetKnownAccs);
	r.GET("/getanticheatlogs", HttpReqGetAntiCheatLogs);

	r.GET("/overridevpn", HttpReqOverrideVPN);

	r.POST("/ticketcreate", HttpReqTicketCreate);
	r.POST("/ticketreply", HttpReqTicketReply);
	r.GET("/ticketlist", HttpReqTicketList);
	r.GET("/ticketmessages", HttpReqTicketMessages);

	r.GET("/status", HttpReqStatus);
	r.GET("/validateprofile", HttpReqValidateProf);
	r.GET("/acceptrules", HttpReqAcceptRules);
	r.GET("/acceptban", HttpReqAcceptBan);
	r.GET("/updatenickname", HttpReqUpdateNameAvatar);
	r.GET("/updatenameavatar", HttpReqUpdateNameAvatar);
	r.GET("/logout", HttpReqLogout);

	r.GET("/getonlineplayers", HttpReqGetOnlinePlayers);
	r.GET("/getplayers", HttpReqGetPlayers);

	r.GET("/getbanrecords", HttpReqGetBanRecords);

	r.GET("/getqueue", HttpReqGetQueue);
	r.GET("/joinqueue", HttpReqJoinQueue);
	r.GET("/leavequeue", HttpReqLeaveQueue);
	r.GET("/readyup", HttpReqReadyUp);

	r.GET("/openidcallback", HttpReqOpenID);
	r.GET("/myip", HttpReqMyIP);
	r.GET("/auth", HttpReqAuth);
	r.GET("/home", HttpReqHome);
	r.POST("/home", HttpReqHome);
	r.GET("/steamid32to64", HttpSteamID32to64);

	r.GET("/getgame", HttpReqGetGame);

	r.GET("/getgameservers", HttpReqGetGameServers);
	r.GET("/pingsreceiver", HttpReqPingsReceiver);

	r.GET("/sendglobalchat", HttpReqSendGlobalChat);
	r.GET("/getglobalchat", HttpReqGetGlobalChat);

	r.GET("/getstreams", HttpReqGetStreams);
	r.GET("/removestream", HttpReqRemoveStream);
	r.GET("/twitchcallback", HttpTwitchOpenIDCallback);
	r.GET("/twitchauth", HttpTwitchAuth);

	r.GET("/offerduo", HttpOfferDuo);
	r.GET("/acceptduo", HttpAcceptDuo);
	r.GET("/cancelduo", HttpCancelDuo);

	r.GET("/getmaps", HttpGetMaps);
	r.GET("/confirmmaps", HttpConfirmMaps);
	r.GET("/revokemapsconfirm", HttpRevokeMapsConfirm);
	r.GET("/refreshmaps", HttpRefreshMaps);


	r.POST("/gs/getgame", HttpReqGSGetGame);
	r.POST("/gs/fullrup", HttpReqGSFullReadyUp);
	r.POST("/gs/partialrup", HttpReqGSPartialReadyUp);
	r.POST("/gs/gameresults", HttpReqGSGameResults);
	r.POST("/gs/anticheatlogs", HttpReqGSAntiCheatLogs);
	r.POST("/gs/chatlogs", HttpReqGSChatLogs);
	r.POST("/gs/checkban", HttpReqGSCheckBan);

	r.GET("/smurf/list_updated", HttpReqSMURFListUpdated);

	
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
