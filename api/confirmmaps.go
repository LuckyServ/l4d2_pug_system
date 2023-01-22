package api

import (
	"github.com/gin-gonic/gin"
	"../players"
	"../players/auth"
	"../database"
	"time"
	"strings"
)


func HttpConfirmMaps(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			players.MuPlayers.Lock();
			i64CurTime := time.Now().UnixMilli();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.Access <= -2) {
				mapResponse["error"] = "Sorry, you are banned, you gotta wait until it expires";
			} else if (pPlayer.LastLimitedRequest + 3000 > i64CurTime) {
				mapResponse["error"] = "Too many requests";
			} else {
				mapResponse["success"] = true;
				pPlayer.LastLimitedRequest = i64CurTime;
				pPlayer.CustomMapsConfirmed = i64CurTime;
				go database.UpdatePlayer(database.DatabasePlayer{
					SteamID64:				pPlayer.SteamID64,
					NicknameBase64:			pPlayer.NicknameBase64,
					AvatarSmall:			pPlayer.AvatarSmall,
					AvatarBig:				pPlayer.AvatarBig,
					Mmr:					pPlayer.Mmr,
					MmrUncertainty:			pPlayer.MmrUncertainty,
					LastGameResult:			pPlayer.LastGameResult,
					Access:					pPlayer.Access,
					ProfValidated:			pPlayer.ProfValidated,
					RulesAccepted:			pPlayer.RulesAccepted,
					Twitch:					pPlayer.Twitch,
					CustomMapsConfirmed:	pPlayer.CustomMapsConfirmed,
					LastCampaignsPlayed:	strings.Join(pPlayer.LastCampaignsPlayed, "|"),
					});
				players.I64LastPlayerlistUpdate = i64CurTime;
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
