package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../queue"
)


func HttpReqReadyUp(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];

			if (!pPlayer.IsInQueue) {
				mapResponse["error"] = "You are not in queue";
			} else if (pPlayer.IsReadyConfirmed) {
				mapResponse["error"] = "You have already ReadyUpped";
			} else if (!pPlayer.IsReadyUpRequested) {
				mapResponse["error"] = "ReadyUp isnt requested";
			} else if (!pPlayer.IsOnline) {
				mapResponse["error"] = "Somehow you are not Online, try to refresh the page";
			} else if (pPlayer.Access <= -2) {
				mapResponse["error"] = "Sorry, you are banned, you gotta wait until it expires";
			} else {
				queue.ReadyUp(pPlayer);
				mapResponse["success"] = true;
			}
			players.MuPlayers.Unlock();
		}
	}

	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
