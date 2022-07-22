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
	"github.com/buger/jsonparser"
	"../database"
)


func HttpReqValidateProf(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.ProfValidated) {
				players.MuPlayers.Unlock();
				mapResponse["error"] = "Your profile is already validated";
			} else if (pPlayer.LastValidateReq + settings.ProfValidateCooldown > time.Now().UnixMilli()) {
				players.MuPlayers.Unlock();
				mapResponse["error"] = "Too many validation requests. Try again in a minute.";
			} else {
				pPlayer.LastValidateReq = time.Now().UnixMilli();
				players.MuPlayers.Unlock();

				clientSteam := http.Client{
					Timeout: 10 * time.Second,
				}
				respSteam, errSteam := clientSteam.Get("https://api.steampowered.com/ISteamUserStats/GetUserStatsForGame/v0002/?appid=550&key="+settings.SteamApiKey+"&steamid="+oSession.SteamID64);
				if (errSteam != nil) {
					mapResponse["error"] = "Steam servers did not respond. Try again later.";
				} else {
					if (respSteam.StatusCode != 200) {
						mapResponse["error"] = "Couldnt get your game details. Make sure your L4D2 stats is public, and try again in a minute.";
					} else {
						byResBody, errResBody := ioutil.ReadAll(respSteam.Body);
						if (errResBody != nil) {
							mapResponse["error"] = "Couldnt get your game details. Make sure your L4D2 stats is public, and try again in a minute.";
						} else {
							var i64VersusGamesWon, i64VersusGamesLost int64;
							jsonparser.ArrayEach(byResBody, func(valueStats []byte, dataType jsonparser.ValueType, offset int, err error) {
								sStatsName, _ := jsonparser.GetString(valueStats, "name");

								if (sStatsName == "Stat.GamesWon.Versus") {
									i64Buffer, errBuffer := jsonparser.GetInt(valueStats, "value");
									if (errBuffer == nil) {
										i64VersusGamesWon = i64Buffer;
									}
								} else if (sStatsName == "Stat.GamesLost.Versus") {
									i64Buffer, errBuffer := jsonparser.GetInt(valueStats, "value");
									if (errBuffer == nil) {
										i64VersusGamesLost = i64Buffer;
									}
								}

							}, "playerstats", "stats");
							iVersusGamePlayed := int(i64VersusGamesWon + i64VersusGamesLost);
							if (iVersusGamePlayed >= settings.MinVersusGamesPlayed) {
								mapResponse["success"] = true;

								players.MuPlayers.Lock();
								pPlayer.ProfValidated = true;
								pPlayer.LastChanged = time.Now().UnixMilli();
								players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
								iNewMmr := settings.DefaultMaxMmr;
								if (iVersusGamePlayed < settings.DefaultMaxMmr) {
									iNewMmr = iVersusGamePlayed;
								}
								pPlayer.Mmr = iNewMmr;
								go database.UpdatePlayer(database.DatabasePlayer{
									SteamID64:			pPlayer.SteamID64,
									NicknameBase64:		pPlayer.NicknameBase64,
									Mmr:				pPlayer.Mmr,
									MmrUncertainty:		pPlayer.MmrUncertainty,
									Access:				pPlayer.Access,
									ProfValidated:		pPlayer.ProfValidated,
									RulesAccepted:		pPlayer.RulesAccepted,
								});
								players.MuPlayers.Unlock();

							} else {
								mapResponse["error"] = fmt.Sprintf("You dont have enough of Versus playtime on your account. Play at least %d public Versus games from the L4D2 menu, and then try the button again.", settings.MinVersusGamesPlayed);
							}
						}
					}
					respSteam.Body.Close();
				}
			}

		} else {
			mapResponse["error"] = "Please authorize first";
		}
	} else {
		mapResponse["error"] = "Please authorize first";
	}
	
	c.Header("Access-Control-Allow-Origin", "*");
	c.JSON(200, mapResponse);
}
