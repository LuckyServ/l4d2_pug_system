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
			if (!players.MapPlayers[oSession.SteamID64].RulesAccepted) {
				mapResponse["success"] = true;
				players.MapPlayers[oSession.SteamID64].RulesAccepted = true;
				go database.UpdatePlayer(database.DatabasePlayer{
					SteamID64:			players.MapPlayers[oSession.SteamID64].SteamID64,
					NicknameBase64:		players.MapPlayers[oSession.SteamID64].NicknameBase64,
					Mmr:				players.MapPlayers[oSession.SteamID64].Mmr,
					MmrUncertainty:		players.MapPlayers[oSession.SteamID64].MmrUncertainty,
					Access:				players.MapPlayers[oSession.SteamID64].Access,
					ProfValidated:		players.MapPlayers[oSession.SteamID64].ProfValidated,
					RulesAccepted:		players.MapPlayers[oSession.SteamID64].RulesAccepted,
					});
				players.MapPlayers[oSession.SteamID64].LastChanged = time.Now().UnixMilli();
				players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
			} else {
				mapResponse["error"] = 2; //already accepted
			}			
			players.MuPlayers.Unlock();
		} else {
			mapResponse["error"] = 1; //unauthorized
		}
	} else {
		mapResponse["error"] = 1; //unauthorized
	}

	c.Header("Access-Control-Allow-Origin", "*");
	c.JSON(200, mapResponse);
}
