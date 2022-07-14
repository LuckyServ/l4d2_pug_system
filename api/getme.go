package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../settings"
	"fmt"
	"time"
)


func HttpReqGetMe(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			mapResponse["success"] = true;
			mapResponse["steamid64"] = 	oSession.SteamID64;

			players.MuPlayers.Lock();

			mapResponse["nickname_base64"] = 	players.MapPlayers[oSession.SteamID64].NicknameBase64;
			mapResponse["mmr"] = 				players.MapPlayers[oSession.SteamID64].Mmr;
			mapResponse["access"] = 			players.MapPlayers[oSession.SteamID64].Access;
			mapResponse["profile_validated"] = 	players.MapPlayers[oSession.SteamID64].ProfValidated;
			mapResponse["rules_accepted"] = 	players.MapPlayers[oSession.SteamID64].RulesAccepted;
			mapResponse["is_online"] = 			players.MapPlayers[oSession.SteamID64].IsOnline;
			mapResponse["is_ingame"] = 			players.MapPlayers[oSession.SteamID64].IsInGame;
			mapResponse["is_inlobby"] = 		players.MapPlayers[oSession.SteamID64].IsInLobby;

			if (players.MapPlayers[oSession.SteamID64].MmrUncertainty <= settings.MmrStable) {
				mapResponse["mmr_certain"] = true;
			} else {
				mapResponse["mmr_certain"] = false;
			}

			players.MuPlayers.Unlock();
		}
	}
	
	c.Header("Access-Control-Allow-Origin", "*");
	c.SetCookie("player_updated_at", fmt.Sprintf("%d", time.Now().UnixMilli()), 2592000, "/", "", true, false);
	c.JSON(200, mapResponse);
}
