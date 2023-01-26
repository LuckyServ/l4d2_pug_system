package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../bans"
	"strconv"
)


func HttpReqDeleteBan(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	i64CreatedAt, errCreatedAt := strconv.ParseInt(c.Query("created_at"), 10, 64);

	mapResponse["success"] = false;
	if (errCreatedAt == nil && i64CreatedAt > 0) {
		if (errCookieSessID == nil && sCookieSessID != "") {
			oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
			if (bAuthorized) {
				players.MuPlayers.RLock();
				iAccess := players.MapPlayers[oSession.SteamID64].Access;
				players.MuPlayers.RUnlock();
				if (iAccess == 4) { //only admin can delete ban

					bans.ChanDeleteBan <- i64CreatedAt;

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
