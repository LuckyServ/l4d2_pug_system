package main

import (
	"github.com/gin-gonic/gin"
	"./settings"
	"./players/auth"
	"io/ioutil"
	"fmt"
	"time"
	"./players"
	"./database"
    "strconv"
	"github.com/yohcop/openid-go"
	"regexp"
	"github.com/antchfx/xmlquery"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"github.com/buger/jsonparser"
)

type PlayerResponse struct {
	SteamID64		string		`json:"steamid64"`
	NicknameBase64	string		`json:"nickname_base64"`
	Mmr				int			`json:"mmr"`
	Access			int 		`json:"access"` //-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	IsInGame		bool		`json:"is_ingame"`
	IsInLobby		bool		`json:"is_inlobby"`
	MmrCertain		bool		`json:"mmr_certain"`
}

type NoOpDiscoveryCache struct{};
var nonceStore = openid.NewSimpleNonceStore();
var discoveryCache = &NoOpDiscoveryCache{};
func (n *NoOpDiscoveryCache) Put(id string, info openid.DiscoveredInfo) {}
func (n *NoOpDiscoveryCache) Get(id string) openid.DiscoveredInfo {
	return nil;
}


func ginInit() {
	gin.SetMode(gin.ReleaseMode); //disable debug logs
	gin.DefaultWriter = ioutil.Discard; //disable output
	r := gin.Default();
	r.MaxMultipartMemory = 1 << 20;

	r.GET("/status", HttpReqStatus);
	r.POST("/shutdown", HttpReqShutdown);
	r.GET("/getme", HttpReqGetMe);
	r.GET("/getonlineplayers", HttpReqGetOnlinePlayers);
	r.GET("/openidcallback", HttpReqOpenID);
	r.GET("/validateprofile", HttpValidateProf);
	
	fmt.Printf("Starting web server\n");
	go r.Run(":"+settings.ListenPort); //Listen on port
}



func HttpReqStatus(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	i64CurTime := time.Now().UnixMilli();

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	sCookiePlayersUpdatedAt, _ := c.Cookie("players_updated_at");
	i64CookiePlayersUpdatedAt, _ := strconv.ParseInt(sCookiePlayersUpdatedAt, 10, 64);
	sCookiePlayerUpdatedAt, _ := c.Cookie("player_updated_at");
	i64CookiePlayerUpdatedAt, _ := strconv.ParseInt(sCookiePlayerUpdatedAt, 10, 64);

	mapResponse["success"] = true;
	mapResponse["no_new_lobbies"] = settings.NoNewLobbies;
	mapResponse["brokenmode"] = settings.BrokenMode;
	mapResponse["time"] = i64CurTime;
	if (i64CookiePlayersUpdatedAt <= players.I64LastPlayerlistUpdate) {
		mapResponse["need_update_players"] = true;
	} else {
		mapResponse["need_update_players"] = false;
	}

	mapResponse["authorized"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			mapResponse["authorized"] = true;
			players.MuPlayers.Lock();
			players.UpdatePlayerActivity(oSession.SteamID64);
			if (i64CookiePlayerUpdatedAt <= players.MapPlayers[oSession.SteamID64].LastChanged) {
				mapResponse["need_update_player"] = true;
			} else {
				mapResponse["need_update_player"] = false;
			}
			players.MuPlayers.Unlock();
		}
	}

	c.Header("Access-Control-Allow-Origin", "*");
	c.JSON(200, mapResponse);
}

func HttpReqGetMe(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			mapResponse["success"] = true;
			mapResponse["steamid64"] = 	oSession.SteamID64;

			players.MuPlayers.Lock();

			mapResponse["nickname_base64"] = 	players.MapPlayers[oSession.SteamID64].NicknameBase64;
			mapResponse["mmr"] = 				players.MapPlayers[oSession.SteamID64].Mmr;
			mapResponse["access"] = 			players.MapPlayers[oSession.SteamID64].Access;
			mapResponse["profile_validated"] = 	players.MapPlayers[oSession.SteamID64].ProfValidated;
			mapResponse["rules_accepted"] = 	players.MapPlayers[oSession.SteamID64].RulesAccepted;
			mapResponse["is_online"] = 			players.MapPlayers[oSession.SteamID64].IsOnline;
			mapResponse["is_ingame"] = 			players.MapPlayers[oSession.SteamID64].IsInGame;
			mapResponse["is_inlobby"] = 		players.MapPlayers[oSession.SteamID64].IsInLobby;

			if (players.MapPlayers[oSession.SteamID64].MmrUncertainty <= settings.MmrStable) {
				mapResponse["mmr_certain"] = true;
			} else {
				mapResponse["mmr_certain"] = false;
			}

			players.MuPlayers.Unlock();
		}
	}
	
	c.Header("Access-Control-Allow-Origin", "*");
	c.SetCookie("player_updated_at", fmt.Sprintf("%d", time.Now().UnixMilli()), 2592000, "/", "", true, false);
	c.JSON(200, mapResponse);
}

func HttpReqGetOnlinePlayers(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	var arPlayers []PlayerResponse;
	var iActiveCount, iOnlineCount, iInLobbyCount, iInGameCount int;

	players.MuPlayers.Lock();
	for _, oPlayer := range players.ArrayPlayers {
		if ((oPlayer.IsOnline || oPlayer.IsInGame || oPlayer.IsInLobby) && oPlayer.ProfValidated && oPlayer.RulesAccepted && oPlayer.Access >= -1/*not banned*/) {
			arPlayers = append(arPlayers, PlayerResponse{
				SteamID64:		oPlayer.SteamID64,
				NicknameBase64:	oPlayer.NicknameBase64,
				Mmr:			oPlayer.Mmr,
				Access:			oPlayer.Access,
				IsInGame:		oPlayer.IsInGame,
				IsInLobby:		oPlayer.IsInLobby,
			});
			if (oPlayer.IsInGame) {
				iInGameCount++;
			} else if (oPlayer.IsInLobby) {
				iInLobbyCount++;
			} else if (oPlayer.IsOnline) {
				iOnlineCount++;
			}
		}
	}
	players.MuPlayers.Unlock();
	iActiveCount = iOnlineCount + iInLobbyCount + iInGameCount;

	//sort
	iSize := len(arPlayers);
	if (iSize > 1) {
		bSorted := false;
		for !bSorted {
			bSorted = true;
			for i := 1; i < iSize; i++ {
				if (arPlayers[i].Mmr > arPlayers[i - 1].Mmr) {
					arPlayers[i], arPlayers[i - 1] = arPlayers[i - 1], arPlayers[i]; //switch
					bSorted = false;
				}
			}
		}
	}


	mapResponse["success"] = true;
	mapResponse["count"] = map[string]int{"online": iActiveCount, "in_lobby": iInLobbyCount, "in_game": iInGameCount};
	mapResponse["list"] = arPlayers;

	
	c.Header("Access-Control-Allow-Origin", "*");
	c.SetCookie("players_updated_at", fmt.Sprintf("%d", time.Now().UnixMilli()), 2592000, "/", "", true, false);
	c.JSON(200, mapResponse);
}

func HttpReqShutdown(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	if (!auth.Backend(c.PostForm("backend_auth"))) {
		mapResponse["success"] = false;
		c.JSON(200, mapResponse);
		return;
	}

	mapResponse["success"] = true;

	c.Header("Access-Control-Allow-Origin", "*");
	c.JSON(200, mapResponse);
	go PerformShutDown();
}

func HttpReqOpenID(c *gin.Context) {
	mapParameters := c.Request.URL.Query();

	//Check if Steam url valid
	if _, ok := mapParameters["openid.op_endpoint"]; !ok {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	if (len(mapParameters["openid.op_endpoint"]) <= 0) {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	if (mapParameters["openid.op_endpoint"][0] != "https://steamcommunity.com/openid/login") {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}

	//Validate auth request with Steam
	sReqString := "?dummy=1";
	for sKey, arValues := range mapParameters {
		if (len(arValues) > 0 && sKey != "openid.mode") {
			sReqString = fmt.Sprintf("%s&%s=%s", sReqString, sKey, url.QueryEscape(arValues[0]));
		}
	}
	fullURL := "https://"+settings.BackendDomain + c.Request.URL.Path + sReqString;
	id, err := openid.Verify(fullURL, discoveryCache, nonceStore);
	if (err != nil) {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	vRegEx := regexp.MustCompile(`[0-9]{17}`);
	bySteamID64 := vRegEx.Find([]byte(id));
	if (bySteamID64 == nil) {
		c.Redirect(303, "https://"+settings.HomeDomain);
	}

	//Here is authorized SteamID64
	sSteamID64 := string(bySteamID64);

	//Get nickname
	sNickname := "unknown";
	clientSteam := http.Client{
		Timeout: 15 * time.Second,
	}
	respSteam, errSteam := clientSteam.Get("https://steamcommunity.com/profiles/"+sSteamID64+"/?xml=1");
	if (errSteam != nil) {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	defer respSteam.Body.Close();
	if (respSteam.StatusCode != 200) {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	byResult, errBody := ioutil.ReadAll(respSteam.Body);
	sResult := string(byResult);
	if (errBody != nil || sResult == "") {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	doc, errXML := xmlquery.Parse(strings.NewReader(sResult));
	if (errXML != nil) {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	root := xmlquery.FindOne(doc, "//profile");
	if n := root.SelectElement("//steamID"); n != nil {
		sNickname = n.InnerText();
	}

	//Add auth to the database
	sSessionID := players.AddPlayerAuth(sSteamID64, base64.StdEncoding.EncodeToString([]byte(sNickname)));

	//Set cookie
	c.SetCookie("session_id", sSessionID, 2592000, "/", "", true, false);

	//Redirect to home page
	c.Header("Access-Control-Allow-Origin", "*");
	c.Redirect(303, "https://"+settings.HomeDomain);
}

func HttpValidateProf(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			players.MuPlayers.Lock();
			if (players.MapPlayers[oSession.SteamID64].ProfValidated) {
				players.MuPlayers.Unlock();
				mapResponse["error"] = 2; //already validated
			} else if (players.MapPlayers[oSession.SteamID64].LastValidateReq + 60000/*60sec*/ > time.Now().UnixMilli()) {
				players.MuPlayers.Unlock();
				mapResponse["error"] = 3; //too many requests, wait
			} else {
				players.MapPlayers[oSession.SteamID64].LastValidateReq = time.Now().UnixMilli();
				players.MuPlayers.Unlock();

				clientSteam := http.Client{
					Timeout: 10 * time.Second,
				}
				respSteam, errSteam := clientSteam.Get("https://api.steampowered.com/ISteamUserStats/GetUserStatsForGame/v0002/?appid=550&key="+settings.SteamApiKey+"&steamid="+oSession.SteamID64);
				if (errSteam != nil) {
					mapResponse["error"] = 4; //Steam request error
				} else {
					if (respSteam.StatusCode != 200) {
						mapResponse["error"] = 4; //Steam request error
					} else {
						byResBody, errResBody := ioutil.ReadAll(respSteam.Body);
						if (errResBody != nil) {
							mapResponse["error"] = 4; //Steam request error
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
								players.MapPlayers[oSession.SteamID64].ProfValidated = true;
								players.MapPlayers[oSession.SteamID64].LastChanged = time.Now().UnixMilli();
								players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
								iNewMmr := settings.DefaultMaxMmr;
								if (iVersusGamePlayed < settings.DefaultMaxMmr) {
									iNewMmr = iVersusGamePlayed;
								}
								players.MapPlayers[oSession.SteamID64].Mmr = iNewMmr;
								go database.UpdatePlayer(database.DatabasePlayer{
									SteamID64:			players.MapPlayers[oSession.SteamID64].SteamID64,
									NicknameBase64:		players.MapPlayers[oSession.SteamID64].NicknameBase64,
									Mmr:				players.MapPlayers[oSession.SteamID64].Mmr,
									MmrUncertainty:		players.MapPlayers[oSession.SteamID64].MmrUncertainty,
									Access:				players.MapPlayers[oSession.SteamID64].Access,
									ProfValidated:		players.MapPlayers[oSession.SteamID64].ProfValidated,
									RulesAccepted:		players.MapPlayers[oSession.SteamID64].RulesAccepted,
								});
								players.MuPlayers.Unlock();

							} else {
								mapResponse["error"] = 5; //Not enough games played (or JSON parsing error)
							}
						}
					}
					respSteam.Body.Close();
				}
			}

		} else {
			mapResponse["error"] = 1; //unauthorized
		}
	} else {
		mapResponse["error"] = 1; //unauthorized
	}
	
	c.Header("Access-Control-Allow-Origin", "*");
	c.JSON(200, mapResponse);
}
