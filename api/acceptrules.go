package api

import (
	"github.com/gin-gonic/gin"
	"../players"
	"../players/auth"
	"../database"
	"time"
)


func HttpReqAcceptRules(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (!pPlayer.RulesAccepted) {
				mapResponse["success"] = true;
				pPlayer.RulesAccepted = true;
				go database.UpdatePlayer(database.DatabasePlayer{
					SteamID64:			pPlayer.SteamID64,
					NicknameBase64:		pPlayer.NicknameBase64,
					Mmr:				pPlayer.Mmr,
					MmrUncertainty:		pPlayer.MmrUncertainty,
					Access:				pPlayer.Access,
					ProfValidated:		pPlayer.ProfValidated,
					RulesAccepted:		pPlayer.RulesAccepted,
					});
				pPlayer.LastChanged = time.Now().UnixMilli();
				players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
			} else {
				mapResponse["error"] = "The rules are already accepted";
			}			
			players.MuPlayers.Unlock();
		} else {
			mapResponse["error"] = "Please authorize first";
		}
	} else {
		mapResponse["error"] = "Please authorize first";
	}

	c.Header("Access-Control-Allow-Origin", "*");
	c.JSON(200, mapResponse);
}
