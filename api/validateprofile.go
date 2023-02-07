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
	"../smurf"
	"strings"
	"encoding/base64"
)

var sProfileClosed string = "Couldnt get your game details. Make sure your L4D2 stats is public, and try again in a minute. If you have just made your L4D2 stats public, you have to wait a few minutes before its available via api.";


func HttpReqValidateProf(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			players.MuPlayers.RLock();
			i64CurTime := time.Now().UnixMilli();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.ProfValidated) {
				players.MuPlayers.RUnlock();
				mapResponse["error"] = "Your profile is already validated";
			} else if (pPlayer.LastExternalRequest + settings.ExternalAPICooldown > i64CurTime) {
				players.MuPlayers.RUnlock();
				mapResponse["error"] = fmt.Sprintf("Too many validation requests. Try again in %d seconds.", ((pPlayer.LastExternalRequest + settings.ExternalAPICooldown) - i64CurTime) / 1000);
			} else {
				sCookieUniqueKey, _ := c.Cookie("auth2");
				byNickname, _ := base64.StdEncoding.DecodeString(pPlayer.NicknameBase64);
				go smurf.AnnounceIPAndKey(pPlayer.SteamID64, c.ClientIP(), string(byNickname), sCookieUniqueKey);
				players.MuPlayers.RUnlock();

				players.MuPlayers.Lock();
				pPlayer.LastExternalRequest = time.Now().UnixMilli();
				players.MuPlayers.Unlock();

				clientSteam := http.Client{
					Timeout: 10 * time.Second,
				}
				respSteam, errSteam := clientSteam.Get("https://api.steampowered.com/ISteamUserStats/GetUserStatsForGame/v0002/?appid=550&key="+settings.SteamApiKey+"&steamid="+oSession.SteamID64);
				if (errSteam != nil) {
					mapResponse["error"] = "Steam servers did not respond. Try again later.";
				} else {
					if (respSteam.StatusCode != 200) {
						mapResponse["error"] = sProfileClosed;
					} else {
						byResBody, errResBody := ioutil.ReadAll(respSteam.Body);
						if (errResBody != nil) {
							mapResponse["error"] = sProfileClosed;
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
							if (iVersusGamePlayed >= settings.MinVersusGamesPlayed && iVersusGamePlayed < 999999) {
								sClientIP := c.ClientIP();
								if (smurf.IsVPN(sClientIP)) {
									mapResponse["error"] = "We believe you are using VPN/proxy connection. You can not validate your profile via VPN/proxy. Make sure you disable VPN/proxy and then try again. If you are not using VPN/proxy, and the system falsely accuses you of using it, please contact admin.";
								} else if (!smurf.IsNotVPN(sClientIP)) {
									mapResponse["error"] = "We could not run some checks on your profile. Try again in a minute, and if it still fails, please contact admin.";
								} else {

									mapResponse["success"] = true;

									iMmrShift := database.GetMmrShift();

									players.MuPlayers.Lock();
									pPlayer.ProfValidated = true;
									players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
									iNewMmr := settings.DefaultMaxMmr;
									pPlayer.Mmr = iNewMmr + iMmrShift;
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
									go database.UpdateInitialGames(database.DatabasePlayer{
										SteamID64:			pPlayer.SteamID64,
										InitialGames:		iVersusGamePlayed,
										});
									players.MuPlayers.Unlock();

								}

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
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
