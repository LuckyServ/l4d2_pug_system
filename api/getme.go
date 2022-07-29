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
			mapResponse["steamid64"] = oSession.SteamID64;

			players.MuPlayers.Lock();

			pPlayer := players.MapPlayers[oSession.SteamID64];

			mapResponse["nickname_base64"] = 	pPlayer.NicknameBase64;
			mapResponse["mmr"] = 				pPlayer.Mmr;
			mapResponse["access"] = 			pPlayer.Access;
			mapResponse["profile_validated"] = 	pPlayer.ProfValidated;
			mapResponse["rules_accepted"] = 	pPlayer.RulesAccepted;
			mapResponse["is_online"] = 			pPlayer.IsOnline;
			mapResponse["is_ingame"] = 			pPlayer.IsInGame;
			mapResponse["is_inlobby"] = 		pPlayer.IsInLobby;
			mapResponse["is_idle"] = 			pPlayer.IsIdle;

			if (pPlayer.MmrUncertainty <= settings.MmrStable) {
				mapResponse["mmr_certain"] = true;
			} else {
				mapResponse["mmr_certain"] = false;
			}

			players.MuPlayers.Unlock();
		}
	}
	
	c.Header("Access-Control-Allow-Origin", "https://"+settings.HomeDomain);
	c.Header("Access-Control-Allow-Credentials", "true");
	c.SetCookie("player_updated_at", fmt.Sprintf("%d", time.Now().UnixMilli()), 2592000, "/", "", true, false);
	c.JSON(200, mapResponse);
}
