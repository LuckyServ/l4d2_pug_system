package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../lobby"
)


func HttpReqLeaveLobby(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (!pPlayer.IsInLobby) {
				mapResponse["error"] = 2; //isn't in lobby
			} else if (!pPlayer.IsOnline) {
				mapResponse["error"] = 3; //not online, wtf bro?
			} else {
				//Leave lobby
				if (lobby.Leave(pPlayer)) {
					mapResponse["success"] = true;
				} else {
					mapResponse["error"] = 4; //???
				}
			}
			players.MuPlayers.Unlock();
		} else {
			mapResponse["error"] = 1; //unauthorized
		}
	} else {
		mapResponse["error"] = 1; //unauthorized
	}
	
	c.Header("Access-Control-Allow-Origin", "*");
	c.JSON(200, mapResponse);
}
