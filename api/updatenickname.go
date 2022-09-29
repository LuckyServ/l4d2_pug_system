package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"time"
	"net/http"
	"io/ioutil"
	"../settings"
	"../database"
	"encoding/base64"
	"github.com/buger/jsonparser"
)

var sErr string = "Error retrieving nickname";

func HttpReqUpdateNickname(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			players.MuPlayers.RLock();
			i64CurTime := time.Now().UnixMilli();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.Access <= -2) {
				players.MuPlayers.RUnlock();
				mapResponse["error"] = "Sorry, you are banned, you have to wait until it expires";
			} else if (pPlayer.LastSteamRequest + settings.SteamAPICooldown > i64CurTime) {
				players.MuPlayers.RUnlock();
				mapResponse["error"] = fmt.Sprintf("You cant request name update that often. Try again in %d seconds.", ((pPlayer.LastSteamRequest + settings.SteamAPICooldown) - i64CurTime) / 1000);
			} else if (!pPlayer.ProfValidated) {
				players.MuPlayers.RUnlock();
				mapResponse["error"] = "Please validate your profile first";
			} else if (!pPlayer.RulesAccepted) {
				players.MuPlayers.RUnlock();
				mapResponse["error"] = "Please accept our rules first";
			} else if (!pPlayer.IsOnline) {
				players.MuPlayers.RUnlock();
				mapResponse["error"] = "Somehow you are not Online, try to refresh the page";
			} else {
				players.MuPlayers.RUnlock();
				players.MuPlayers.Lock();
				pPlayer.LastSteamRequest = i64CurTime;
				players.MuPlayers.Unlock();
				clientSteam := http.Client{
					Timeout: 10 * time.Second,
				}
				respSteam, errSteam := clientSteam.Get("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key="+settings.SteamApiKey+"&steamids="+oSession.SteamID64);
				if (errSteam == nil) {
					if (respSteam.StatusCode == 200) {
						byResult, _ := ioutil.ReadAll(respSteam.Body);
						sName, errName := jsonparser.GetString(byResult, "response", "players", "[0]", "personaname");
						sSteamID64, _ := jsonparser.GetString(byResult, "response", "players", "[0]", "steamid");
						if (sSteamID64 == oSession.SteamID64 && errName == nil && sName != "") {

							mapResponse["success"] = true;
							players.MuPlayers.Lock();
							pPlayer.NicknameBase64 = base64.StdEncoding.EncodeToString([]byte(sName));
							players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
							go database.UpdatePlayer(database.DatabasePlayer{
								SteamID64:			pPlayer.SteamID64,
								NicknameBase64:		pPlayer.NicknameBase64,
								Mmr:				pPlayer.Mmr,
								MmrUncertainty:		pPlayer.MmrUncertainty,
								LastGameResult:		pPlayer.LastGameResult,
								Access:				pPlayer.Access,
								ProfValidated:		pPlayer.ProfValidated,
								RulesAccepted:		pPlayer.RulesAccepted,
								});
							players.MuPlayers.Unlock();

						} else {mapResponse["error"] = sErr;}
					} else {mapResponse["error"] = sErr;}
					respSteam.Body.Close();
				} else {mapResponse["error"] = sErr;}
			}
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
