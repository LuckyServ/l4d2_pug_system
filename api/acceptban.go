package api

import (
	"github.com/gin-gonic/gin"
	"../players"
	"../players/auth"
	"../database"
	"time"
)


func HttpReqAcceptBan(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];

			if (pPlayer.Access != -2) {
				mapResponse["error"] = "You are not banned";
			} else if (pPlayer.BanAcceptedAt > 0) {
				mapResponse["error"] = "You have already confirmed that you read this notification";
			} else {
				mapResponse["success"] = true;
				pPlayer.BanAcceptedAt = time.Now().UnixMilli();
				go database.UpdatePlayer(database.DatabasePlayer{
					SteamID64:			pPlayer.SteamID64,
					NicknameBase64:		pPlayer.NicknameBase64,
					Mmr:				pPlayer.Mmr,
					MmrUncertainty:		pPlayer.MmrUncertainty,
					Access:				pPlayer.Access,
					ProfValidated:		pPlayer.ProfValidated,
					RulesAccepted:		pPlayer.RulesAccepted,
					});
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
