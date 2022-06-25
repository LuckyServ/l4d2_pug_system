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
)


func ginInit() {
	gin.SetMode(gin.ReleaseMode); //disable debug logs
	gin.DefaultWriter = ioutil.Discard; //disable output
	r := gin.Default();
	r.MaxMultipartMemory = 1 << 20;

	r.POST("/status", HttpReqStatus);
	//r.POST("/shutdown", HttpReqShutdown);
	r.POST("/updateactivity", HttpReqUpdateActivity);
	r.POST("/getme", HttpReqGetMe);
	
	fmt.Printf("Starting web server\n");
	go r.Run(":"+settings.ListenPort); //Listen on port
}



func HttpReqStatus(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	i64CurTime := time.Now().UnixMilli();

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	sCookiePlayersUpdatedAt, _ := c.Cookie("players_updated_at");
	i64CookiePlayersUpdatedAt, _ := strconv.ParseInt(sCookiePlayersUpdatedAt, 10, 64)
	sCookiePlayerUpdatedAt, _ := c.Cookie("player_updated_at");
	i64CookiePlayerUpdatedAt, _ := strconv.ParseInt(sCookiePlayerUpdatedAt, 10, 64)

	mapResponse["success"] = true;
	mapResponse["shutdown"] = bStateShutdown;
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
			if (i64CookiePlayerUpdatedAt <= players.MapPlayers[oSession.SteamID64].LastChanged) {
				mapResponse["need_update_player"] = true;
			} else {
				mapResponse["need_update_player"] = false;
			}
			players.MuPlayers.Unlock();
		}
	}

	c.JSON(200, mapResponse);
}


func HttpReqUpdateActivity(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			players.UpdatePlayerActivity(oSession.SteamID64);
			mapResponse["success"] = true;
		}
	}
	
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

			/*var bMmrCertain, bPingsUpdated bool;
			if (MapPlayers[sSteamID64].MmrUncertainty <= settings.MmrStable) {
				bMmrCertain = true;
			} else {
				bMmrCertain = false;
			}
			if ((time.Now().UnixMilli() - MapPlayers[sSteamID64].LastPingsUpdate) < settings.PingsMaxAge) {
				bPingsUpdated = true;
			} else {
				bPingsUpdated = false;
			}*/

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
			if ((time.Now().UnixMilli() - players.MapPlayers[oSession.SteamID64].LastPingsUpdate) <= settings.PingsMaxAge) {
				mapResponse["need_update_pings"] = false;
			} else {
				mapResponse["need_update_pings"] = true;
			}

			players.MuPlayers.Unlock();
		}
	}
	
	c.JSON(200, mapResponse);
}





//Old bad api

/*func HttpReqShutdown(c *gin.Context) {

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

	c.JSON(200, mapResponse);
	go PerformShutDown();
}

func HttpReqAddAuth(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	if (!auth.Backend(c.PostForm("backend_auth"))) {
		mapResponse["success"] = false;
		mapResponse["error"] = 1;
		c.JSON(200, mapResponse);
		return;
	}


	sSteamID64 := c.PostForm("steamid64");
	sNicknameBase64 := c.PostForm("nickname_base64");
	if (sSteamID64 == "" || sNicknameBase64 == "") {
		mapResponse["success"] = false;
		mapResponse["error"] = 2;
		c.JSON(200, mapResponse);
		return;
	}

	sSessionID := auth.AddPlayerAuth(sSteamID64, sNicknameBase64);
	mapResponse["success"] = true;
	mapResponse["session_id"] = sSessionID;
	
	c.JSON(200, mapResponse);
}

func HttpReqRemoveAuth(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	if (!auth.Backend(c.PostForm("backend_auth"))) {
		mapResponse["success"] = false;
		mapResponse["error"] = 1;
		c.JSON(200, mapResponse);
		return;
	}


	sSessID := c.PostForm("session_id");
	if (sSessID == "") {
		mapResponse["success"] = false;
		mapResponse["error"] = 2;
		c.JSON(200, mapResponse);
		return;
	}

	bSuccess, iError := auth.RemovePlayerAuth(sSessID);
	mapResponse["success"] = bSuccess;
	if (!bSuccess) {
		mapResponse["error"] = iError;
	}
	
	c.JSON(200, mapResponse);
}

func HttpReqUpdateActivity(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	if (!auth.Backend(c.PostForm("backend_auth"))) {
		mapResponse["success"] = false;
		mapResponse["error"] = 1;
		c.JSON(200, mapResponse);
		return;
	}


	sSteamID64 := c.PostForm("steamid64");
	if (sSteamID64 == "") {
		mapResponse["success"] = false;
		mapResponse["error"] = 2;
		c.JSON(200, mapResponse);
		return;
	}

	bSuccess, iError := players.UpdatePlayerActivity(sSteamID64);
	mapResponse["success"] = bSuccess;
	if (!bSuccess) {
		mapResponse["error"] = iError;
	}
	
	c.JSON(200, mapResponse);
}
*/