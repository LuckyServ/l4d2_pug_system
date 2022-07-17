package api

import (
	"github.com/gin-gonic/gin"
	"time"
    "strconv"
	"../settings"
	"../players"
	"../lobby"
	"../players/auth"
)


func HttpReqStatus(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	i64CurTime := time.Now().UnixMilli();

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	sCookiePlayersUpdatedAt, _ := c.Cookie("players_updated_at");
	i64CookiePlayersUpdatedAt, _ := strconv.ParseInt(sCookiePlayersUpdatedAt, 10, 64);
	sCookiePlayerUpdatedAt, _ := c.Cookie("player_updated_at");
	i64CookiePlayerUpdatedAt, _ := strconv.ParseInt(sCookiePlayerUpdatedAt, 10, 64);
	sCookieLobbiesUpdatedAt, _ := c.Cookie("lobbies_updated_at");
	i64CookieLobbiesUpdatedAt, _ := strconv.ParseInt(sCookieLobbiesUpdatedAt, 10, 64);
	sWindowActive := c.Query("active");

	mapResponse["success"] = true;
	mapResponse["no_new_lobbies"] = settings.NoNewLobbies;
	mapResponse["brokenmode"] = settings.BrokenMode;
	mapResponse["time"] = i64CurTime;

	mapResponse["authorized"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			mapResponse["authorized"] = true;
			players.MuPlayers.Lock();

			if (sWindowActive == "true") {
				players.UpdatePlayerActivity(oSession.SteamID64);
			}

			if (i64CookiePlayerUpdatedAt <= players.MapPlayers[oSession.SteamID64].LastChanged) {
				mapResponse["need_update_player"] = true;
			} else {
				mapResponse["need_update_player"] = false;
			}

			players.MuPlayers.Unlock();
		}
	}

	if (i64CookiePlayersUpdatedAt <= players.I64LastPlayerlistUpdate) {
		mapResponse["need_update_players"] = true;
	} else {
		mapResponse["need_update_players"] = false;
	}
	if (i64CookieLobbiesUpdatedAt <= lobby.I64LastLobbyListUpdate) {
		mapResponse["need_update_lobbies"] = true;
	} else {
		mapResponse["need_update_lobbies"] = false;
	}

	c.Header("Access-Control-Allow-Origin", "*");
	c.JSON(200, mapResponse);
}
