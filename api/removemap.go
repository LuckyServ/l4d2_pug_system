package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../games"
	"../settings"
)


func HttpRemoveMap(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	sCampaignName := c.Query("name");

	mapResponse["success"] = false;
	if (sCampaignName != "") {
		if (errCookieSessID == nil && sCookieSessID != "") {
			oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
			if (bAuthorized) {
				players.MuPlayers.RLock();
				iAccess := players.MapPlayers[oSession.SteamID64].Access;
				players.MuPlayers.RUnlock();
				if (iAccess == 4) { //only admin can remove maps

					var iIndex int = -1;
					games.MuGames.Lock();
					for i, oCampaign:= range settings.MapPool {
						if (oCampaign.Name == sCampaignName) {
							iIndex = i;
						}
					}
					if (iIndex >= 0) {
						mapResponse["success"] = true;

						settings.MapPool = append(settings.MapPool[:iIndex], settings.MapPool[iIndex+1:]...);

					} else {
						mapResponse["error"] = "Map not found";
					}
					games.MuGames.Unlock();

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
