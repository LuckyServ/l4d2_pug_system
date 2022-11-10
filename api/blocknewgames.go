package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../lobby"
)


func HttpReqBlockNewGames(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			players.MuPlayers.RLock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.Access == 4) { //admin
				players.MuPlayers.RUnlock();

				mapResponse["success"] = true;
				lobby.NewLobbiesBlocked = true;

				players.MuPlayers.Lock();
				lobby.MuLobbies.Lock();

				for _, pPlayer := range players.ArrayPlayers {
					lobby.Leave(pPlayer, true);
					pPlayer.IsAutoSearching = false;
				}

				lobby.MuLobbies.Unlock();
				players.MuPlayers.Unlock();

			} else {
				players.MuPlayers.RUnlock();
				mapResponse["error"] = "You dont have access to this command";
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
