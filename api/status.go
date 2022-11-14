package api

import (
	"github.com/gin-gonic/gin"
	"time"
    "strconv"
	"../settings"
	"../players"
	"../chat"
	"../smurf"
	"../utils"
	"../players/auth"
	"regexp"
	"../queue"
)


func HttpReqStatus(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	i64CurTime := time.Now().UnixMilli();

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	sCookiePlayersUpdatedAt, _ := c.Cookie("players_updated_at");
	i64CookiePlayersUpdatedAt, _ := strconv.ParseInt(sCookiePlayersUpdatedAt, 10, 64);
	sCookieQueueUpdatedAt, _ := c.Cookie("queue_updated_at");
	i64CookieQueueUpdatedAt, _ := strconv.ParseInt(sCookieQueueUpdatedAt, 10, 64);
	sCookieGameUpdatedAt, _ := c.Cookie("game_updated_at");
	i64CookieGameUpdatedAt, _ := strconv.ParseInt(sCookieGameUpdatedAt, 10, 64);
	sCookieGlobalChatUpdatedAt, _ := c.Cookie("globalchat_updated_at");
	i64CookieGlobalChatUpdatedAt, _ := strconv.ParseInt(sCookieGlobalChatUpdatedAt, 10, 64);

	mapResponse["success"] = true;
	mapResponse["no_new_games"] = queue.NewGamesBlocked;
	mapResponse["brokenmode"] = settings.BrokenMode;
	mapResponse["time"] = i64CurTime;

	mapResponse["authorized"] = false;
	mapResponse["need_update_game"] = false;
	mapResponse["need_update_queue"] = false;
	mapResponse["need_emit_readyup_sound"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			mapResponse["authorized"] = true;
			go smurf.AnnounceIP(c.ClientIP()); //for faster VPN info retrieve in the future
			sCookieUniqueKey, _ := c.Cookie("auth2");
			players.MuPlayers.Lock();
			players.UpdatePlayerActivity(oSession.SteamID64, sCookieUniqueKey, c.ClientIP());
			players.MuPlayers.Unlock();

			players.MuPlayers.RLock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (i64CookieGameUpdatedAt <= pPlayer.LastGameChanged) {
				mapResponse["need_update_game"] = true;
			}
			if (i64CookieQueueUpdatedAt <= pPlayer.LastQueueChanged) {
				mapResponse["need_update_queue"] = true;
			}
			if (pPlayer.IsInQueue && pPlayer.IsReadyUpRequested && !pPlayer.IsReadyConfirmed) {
				mapResponse["need_emit_readyup_sound"] = true;
			}
			players.MuPlayers.RUnlock();
		}
	}

	if (i64CookiePlayersUpdatedAt <= players.I64LastPlayerlistUpdate) {
		mapResponse["need_update_players"] = true;
	} else {
		mapResponse["need_update_players"] = false;
	}
	if (i64CookieGlobalChatUpdatedAt <= chat.I64LastGlobalChatUpdate) {
		mapResponse["need_update_globalchat"] = true;
	} else {
		mapResponse["need_update_globalchat"] = false;
	}

	sCookieUniqueKey, _ := c.Cookie("auth2");
	bKeyValid, _ := regexp.MatchString(`^[0-9a-z]{16,200}$`, sCookieUniqueKey);
	if (!bKeyValid) {
		sRand, _ := utils.GenerateRandomString(40, "0123456789abcdefghijklmnopqrstuvwxyz");
		c.SetCookie("auth2", sRand, 2000000000, "/", "", true, false);
	}
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
