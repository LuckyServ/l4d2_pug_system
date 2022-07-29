package api

import (
	"github.com/gin-gonic/gin"
	"../lobby"
	"../players/auth"
	"../players"
	"../settings"
)


func HttpReqReadyUp(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];

			if (!pPlayer.IsInLobby) {
				mapResponse["error"] = "You are not in lobby";
			} else if (pPlayer.IsReadyInLobby) {
				mapResponse["error"] = "You have already ReadyUpped";
			} else if (!pPlayer.IsOnline) {
				mapResponse["error"] = "Somehow you are not Online, try to refresh the page";
			} else if (pPlayer.Access == -2) {
				mapResponse["error"] = "Sorry, you are banned, you gotta wait until it expires";
			} else {
				lobby.MuLobbies.Lock();
				if (lobby.Ready(pPlayer)) {
					mapResponse["success"] = true;
				} else {
					mapResponse["error"] = "Unknown error, try again";
				}
				lobby.MuLobbies.Unlock();
			}
			players.MuPlayers.Unlock();
		}
	}

	
	c.Header("Access-Control-Allow-Origin", "https://"+settings.HomeDomain);
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
