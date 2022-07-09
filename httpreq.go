package main

import (
	"github.com/gin-gonic/gin"
	"./settings"
	"./players/auth"
	"io/ioutil"
	"fmt"
	"time"
	"./players"
    "strconv"
	"github.com/yohcop/openid-go"
	"regexp"
	"github.com/antchfx/xmlquery"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
)

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
	r.GET("/openidcallback", HttpReqOpenID);
	
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
	mapResponse["shutdown"] = bStateShutdown;
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
			mapResponse["mmr_uncertainty"] = 	players.MapPlayers[oSession.SteamID64].MmrUncertainty;
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

func HttpReqShutdown(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	if (!auth.Backend(c.PostForm("backend_auth"))) {
		mapResponse["success"] = false;
		mapResponse["error"] = 1;
		c.JSON(200, mapResponse);
		return;
	}

	bSetShutdown, iError := SetShutDown();
	if (!bSetShutdown) {
		mapResponse["success"] = false;
		mapResponse["error"] = iError;
	} else {
		mapResponse["success"] = true;
	}

	c.Header("Access-Control-Allow-Origin", "*");
	c.JSON(200, mapResponse);
	go PerformShutDown();
}

func HttpReqOpenID(c *gin.Context) {
	arParameters := c.Request.URL.Query();

	//Check if Steam url valid
	if _, ok := arParameters["openid.op_endpoint"]; !ok {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	if (len(arParameters["openid.op_endpoint"]) <= 0) {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	if (arParameters["openid.op_endpoint"][0] != "https://steamcommunity.com/openid/login") {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}

	//Validate auth request with Steam
	sReqString := "?dummy=1";
	for sKey, arValues := range arParameters {
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
