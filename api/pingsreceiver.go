package api

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"../games"
	"../players/auth"
	"../players"
	"../settings"
	"strconv"
)


func HttpReqPingsReceiver(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.IsInGame) {
				games.MuGames.Lock();
				pGame := games.MapGames[pPlayer.GameID];
				if (pGame.State == games.StateWaitPings) {
					mapResponse["success"] = true;

					for i, _ := range settings.HardwareServers {
						sPingMS := c.Query(settings.HardwareDomains[i]);
						if (sPingMS != "") {
							iPingMS, errPingMS := strconv.Atoi(sPingMS);
							if (errPingMS == nil && iPingMS > 0) {
								iOldPing, bAlrPinged := pPlayer.GameServerPings[settings.HardwareServers[i]];
								if (bAlrPinged) {
									if (iPingMS < iOldPing) {
										pPlayer.GameServerPings[settings.HardwareServers[i]] = iPingMS;
									}
								} else {
									pPlayer.GameServerPings[settings.HardwareServers[i]] = iPingMS;
								}
							}
						}
					}

				} else {
					mapResponse["error"] = "The game isnt expecting a ping info from you";
				}
				games.MuGames.Unlock();
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

	//Testing
	oSession, _ := auth.GetSession(sCookieSessID);
	pPlayer := players.MapPlayers[oSession.SteamID64];
	mapResponse["success"] = true;
	for i, _ := range settings.HardwareServers {
		sPingMS := c.Query(settings.HardwareDomains[i]);
		if (sPingMS != "") {
			iPingMS, errPingMS := strconv.Atoi(sPingMS);
			if (errPingMS == nil && iPingMS > 0) {
				iOldPing, bAlrPinged := pPlayer.GameServerPings[settings.HardwareServers[i]];
				if (bAlrPinged) {
					if (iPingMS < iOldPing) {
						pPlayer.GameServerPings[settings.HardwareServers[i]] = iPingMS;
					}
				} else {
					pPlayer.GameServerPings[settings.HardwareServers[i]] = iPingMS;
				}
			}
		}
	}
	
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
