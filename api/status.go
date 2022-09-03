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
	sCookieLobbiesUpdatedAt, _ := c.Cookie("lobbies_updated_at");
	i64CookieLobbiesUpdatedAt, _ := strconv.ParseInt(sCookieLobbiesUpdatedAt, 10, 64);
	sCookieGameUpdatedAt, _ := c.Cookie("game_updated_at");
	i64CookieGameUpdatedAt, _ := strconv.ParseInt(sCookieGameUpdatedAt, 10, 64);

	mapResponse["success"] = true;
	mapResponse["no_new_lobbies"] = settings.NoNewLobbies;
	mapResponse["brokenmode"] = settings.BrokenMode;
	mapResponse["time"] = i64CurTime;

	mapResponse["authorized"] = false;
	mapResponse["need_update_game"] = false;
	mapResponse["need_emit_readyup_sound"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			mapResponse["authorized"] = true;
			players.MuPlayers.Lock();
			players.UpdatePlayerActivity(oSession.SteamID64);

			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (i64CookieGameUpdatedAt <= pPlayer.LastGameChanged) {
				mapResponse["need_update_game"] = true;
			}
			if (pPlayer.IsInLobby) {
				pLobby := lobby.MapLobbies[pPlayer.LobbyID];
				if (pLobby.PlayerCount >= 8 && !pPlayer.IsReadyInLobby) {
					mapResponse["need_emit_readyup_sound"] = true;
				}
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

	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
