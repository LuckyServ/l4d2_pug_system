package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../bans"
)


func HttpReqUnban(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	sSteamID64 := c.Query("steamid64");

	mapResponse["success"] = false;
	if (sSteamID64 != "") {
		if (errCookieSessID == nil && sCookieSessID != "") {
			oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
			if (bAuthorized) {
				players.MuPlayers.RLock();
				iAccess := players.MapPlayers[oSession.SteamID64].Access;
				players.MuPlayers.RUnlock();
				if (iAccess >= 1) { //only moderators can unban

					bans.ChanUnbanManual <- bans.EntManualUnbanReq{
						SteamID64:			sSteamID64,
						RequestedBy:		oSession.SteamID64,
					};

					mapResponse["success"] = true;
				} else {
					mapResponse["error"] = "You dont have access to this command";
				}
			} else {
				mapResponse["error"] = "Please authorize first";
			}
		} else {
			mapResponse["error"] = "Please authorize first";
		}
	} else {
		mapResponse["error"] = "Bad parameters";
	}
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
