package api

import (
	"github.com/gin-gonic/gin"
	"../players"
	"../players/auth"
	"../bans"
)


func HttpReqAcceptBan(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			players.MuPlayers.RLock();
			pPlayer := players.MapPlayers[oSession.SteamID64];

			if (pPlayer.Access >= -1) {
				players.MuPlayers.RUnlock();
				mapResponse["error"] = "You are not banned";
			} else if (pPlayer.BanAcceptedAt > 0) {
				players.MuPlayers.RUnlock();
				mapResponse["error"] = "You have already confirmed that you read this notification";
			} else {
				players.MuPlayers.RUnlock();
				mapResponse["success"] = true;
				bans.ChanAcceptBan <- oSession.SteamID64;
			}

		} else {
			mapResponse["error"] = "Please authorize first";
		}
	} else {
		mapResponse["error"] = "Please authorize first";
	}

	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
