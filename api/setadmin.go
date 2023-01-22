package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../database"
	"strconv"
	"time"
	"strings"
)


func HttpReqSetAdmin(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	sSteamID64 := c.Query("steamid64");
	iSetAccess, _ := strconv.Atoi(c.Query("access"));

	mapResponse["success"] = false;
	if (sSteamID64 != "" && iSetAccess >= 0 && iSetAccess <= 4) {
		if (errCookieSessID == nil && sCookieSessID != "") {
			oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
			if (bAuthorized) {
				players.MuPlayers.RLock();
				iAccess := players.MapPlayers[oSession.SteamID64].Access;
				players.MuPlayers.RUnlock();
				if (iAccess == 4) { //only admin can add and remove admins

					players.MuPlayers.Lock();

					pPlayer, bFound := players.MapPlayers[sSteamID64];
					if (bFound) {
						pPlayer.Access = iSetAccess;
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
							players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
						mapResponse["success"] = true;
					} else {
						mapResponse["error"] = "No such player";
					}

					players.MuPlayers.Unlock();

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
