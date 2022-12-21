package api

import (
	"github.com/gin-gonic/gin"
	"../players"
	"../streams"
	"../database"
	"../settings"
	"fmt"
	"time"
	"../players/auth"
	"regexp"
)


func HttpReqAddStream(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	sUsername := c.Query("username");
	bUsernameValid, _ := regexp.MatchString(`^[\w]{2,32}$`, sUsername);

	mapResponse["success"] = false;
	if (!bUsernameValid) {
		mapResponse["error"] = "Invalid username";
	} else {
		if (errCookieSessID == nil && sCookieSessID != "") {
			oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
			if (bAuthorized) {
				players.MuPlayers.RLock();
				i64CurTime := time.Now().UnixMilli();
				pPlayer := players.MapPlayers[oSession.SteamID64];
				if (pPlayer.LastExternalRequest + settings.ExternalAPICooldown > i64CurTime) {
					players.MuPlayers.RUnlock();
					mapResponse["error"] = fmt.Sprintf("You cant request Twitch api that often. Try again in %d seconds.", ((pPlayer.LastExternalRequest + settings.ExternalAPICooldown) - i64CurTime) / 1000);
				} else if (!pPlayer.ProfValidated) {
					players.MuPlayers.RUnlock();
					mapResponse["error"] = "Please validate your profile first";
				} else if (!pPlayer.IsOnline) {
					players.MuPlayers.RUnlock();
					mapResponse["error"] = "Somehow you are not Online, try to refresh the page";
				} else {
					players.MuPlayers.RUnlock();
					players.MuPlayers.Lock();
					pPlayer.LastExternalRequest = i64CurTime;
					players.MuPlayers.Unlock();

					sUserID, errUserID := streams.GetTwitchUserID(sUsername);
					if (errUserID != nil) {
						mapResponse["error"] = errUserID.Error();
					} else {
						mapResponse["success"] = true;
						players.MuPlayers.Lock();
						pPlayer.Twitch = sUserID;
						go database.UpdatePlayer(database.DatabasePlayer{
							SteamID64:			pPlayer.SteamID64,
							NicknameBase64:		pPlayer.NicknameBase64,
							AvatarSmall:		pPlayer.AvatarSmall,
							AvatarBig:			pPlayer.AvatarBig,
							Mmr:				pPlayer.Mmr,
							MmrUncertainty:		pPlayer.MmrUncertainty,
							LastGameResult:		pPlayer.LastGameResult,
							Access:				pPlayer.Access,
							ProfValidated:		pPlayer.ProfValidated,
							RulesAccepted:		pPlayer.RulesAccepted,
							Twitch:				pPlayer.Twitch,
							});
						players.MuPlayers.Unlock();
					}
				}
			} else {
				mapResponse["error"] = "Please authorize first";
			}
		} else {
			mapResponse["error"] = "Please authorize first";
		}
	}
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
