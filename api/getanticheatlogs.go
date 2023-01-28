package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../database"
)


func HttpReqGetAntiCheatLogs(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			players.MuPlayers.RLock();
			iAccess := players.MapPlayers[oSession.SteamID64].Access;
			players.MuPlayers.RUnlock();
			if (iAccess >= 2) { //cheat moderator or admin

				mapResponse["success"] = true;
				mapResponse["logs"] = database.GetAnticheatLogs(); //slow, makes connections to database

			} else {
				mapResponse["error"] = "You dont have access to this information";
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
