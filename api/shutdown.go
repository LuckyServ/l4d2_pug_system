package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
)

var ChShutdown chan bool = make(chan bool);


func HttpReqShutdown(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			players.MuPlayers.RLock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.Access == 4) { //admin
				mapResponse["success"] = true;
			}
			players.MuPlayers.RUnlock();
		}
	}
	

	c.JSON(200, mapResponse);
	if (mapResponse["success"] == true) {
		go PerformShutDown();
	}
}

func PerformShutDown() {
	ChShutdown <- true;
}