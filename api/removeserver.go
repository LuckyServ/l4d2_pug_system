package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../games"
	"../settings"
)


func HttpRemoveServer(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	sDomain := c.Query("domain");

	mapResponse["success"] = false;
	if (sDomain != "") {
		if (errCookieSessID == nil && sCookieSessID != "") {
			oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
			if (bAuthorized) {
				players.MuPlayers.RLock();
				iAccess := players.MapPlayers[oSession.SteamID64].Access;
				players.MuPlayers.RUnlock();

				var sServerAdmin string;
				games.MuGames.RLock();
				for _, oServer := range settings.GameServers {
					if (oServer.Domain == sDomain) {
						sServerAdmin = oServer.Admin;
					}
				}
				games.MuGames.RUnlock();

				if (sServerAdmin == "") {
					mapResponse["error"] = "Server not found";
				} else if (iAccess == 4 || sServerAdmin == oSession.SteamID64) { //only admin of the server or admin of l4d2center can remove servers

					var iIndex int = -1;
					games.MuGames.Lock();
					for i, oServer:= range settings.GameServers {
						if (oServer.Domain == sDomain) {
							iIndex = i;
						}
					}
					if (iIndex >= 0) {
						mapResponse["success"] = true;

						settings.GameServers = append(settings.GameServers[:iIndex], settings.GameServers[iIndex+1:]...);

					} else {
						mapResponse["error"] = "Server not found";
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
