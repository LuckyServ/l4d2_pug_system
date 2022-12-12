package api

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"../games"
	"../players/auth"
	"../players"
	"../settings"
	"../smurf"
	"strconv"
)


func HttpReqPingsReceiver(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.IsInGame) {
				games.MuGames.RLock();
				pGame := games.MapGames[pPlayer.GameID];
				if (pGame.State == games.StateWaitPings) {
					mapResponse["success"] = true;

					if (!smurf.IsVPN(c.ClientIP())) { //silently avoid storing ping data from VPN users
						for _, oGameServer := range settings.GameServers {
							sPingMS := c.Query(oGameServer.Domain);
							if (sPingMS != "") {
								iPingMS, errPingMS := strconv.Atoi(sPingMS);
								if (errPingMS == nil && iPingMS > 0) {
									iOldPing, bAlrPinged := pPlayer.GameServerPings[oGameServer.IP];
									if (bAlrPinged) {
										if (iPingMS < iOldPing) {
											pPlayer.GameServerPings[oGameServer.IP] = iPingMS;
										}
									} else {
										pPlayer.GameServerPings[oGameServer.IP] = iPingMS;
									}
								}
							}
						}
					}

				} else {
					mapResponse["error"] = "The game isnt expecting a ping info from you";
				}
				games.MuGames.RUnlock();
			} else {
				mapResponse["error"] = "You are not in game";
			}
			players.MuPlayers.Unlock();
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
