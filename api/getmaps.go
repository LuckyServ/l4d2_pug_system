package api

import (
	"github.com/gin-gonic/gin"
	"../players"
	"../players/auth"
	"../settings"
)


func HttpGetMaps(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = true;
	mapResponse["authorized"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			mapResponse["authorized"] = true;
			players.MuPlayers.RLock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			mapResponse["maps_confirmed_at"] = pPlayer.CustomMapsConfirmed;
			players.MuPlayers.RUnlock();
		}
	}

	mapResponse["campaigns"] = settings.MapPool;
	mapResponse["newest_map"] = settings.NewestCustomMap;

	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
