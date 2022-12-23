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
	"../streams"
)


func HttpReqStatus(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	i64CurTime := time.Now().UnixMilli();

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	i64QueryPlayersUpdatedAt, _ := strconv.ParseInt(c.Query("players_updated_at"), 10, 64);
	i64QueryQueueUpdatedAt, _ := strconv.ParseInt(c.Query("queue_updated_at"), 10, 64);
	i64QueryGameUpdatedAt, _ := strconv.ParseInt(c.Query("game_updated_at"), 10, 64);
	i64QueryGlobalChatUpdatedAt, _ := strconv.ParseInt(c.Query("globalchat_updated_at"), 10, 64);
	i64QueryStreamersUpdatedAt, _ := strconv.ParseInt(c.Query("streamers_updated_at"), 10, 64);

	mapResponse["success"] = true;
	mapResponse["no_new_games"] = queue.NewGamesBlocked;
	mapResponse["brokenmode"] = settings.BrokenMode;
	mapResponse["time"] = i64CurTime;

	mapResponse["authorized"] = false;
	mapResponse["need_update_game"] = false;
	mapResponse["need_update_queue"] = false;
	mapResponse["need_emit_readyup_sound"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			mapResponse["authorized"] = true;
			go smurf.AnnounceIP(c.ClientIP()); //for faster VPN info retrieve in the future
			sCookieUniqueKey, _ := c.Cookie("auth2");
			players.MuPlayers.Lock();
			players.UpdatePlayerActivity(oSession.SteamID64, sCookieUniqueKey, c.ClientIP());
			players.MuPlayers.Unlock();

			players.MuPlayers.RLock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (i64QueryGameUpdatedAt <= pPlayer.LastGameChanged) {
				mapResponse["need_update_game"] = true;
			}
			if (i64QueryQueueUpdatedAt <= pPlayer.LastQueueChanged) {
				mapResponse["need_update_queue"] = true;
			}
			if (pPlayer.IsInQueue && pPlayer.IsReadyUpRequested && !pPlayer.IsReadyConfirmed) {
				mapResponse["need_emit_readyup_sound"] = true;
			}
			players.MuPlayers.RUnlock();
		}
	}

	if (i64QueryPlayersUpdatedAt <= players.I64LastPlayerlistUpdate) {
		mapResponse["need_update_players"] = true;
	} else {
		mapResponse["need_update_players"] = false;
	}
	if (i64QueryGlobalChatUpdatedAt <= chat.I64LastGlobalChatUpdate) {
		mapResponse["need_update_globalchat"] = true;
	} else {
		mapResponse["need_update_globalchat"] = false;
	}
	if (i64QueryStreamersUpdatedAt <= streams.I64LastStreamersUpdate) {
		mapResponse["need_update_streamers"] = true;
	} else {
		mapResponse["need_update_streamers"] = false;
	}

	sCookieUniqueKey, _ := c.Cookie("auth2");
	bKeyValid, _ := regexp.MatchString(`^[0-9a-z]{16,200}$`, sCookieUniqueKey);
	if (!bKeyValid) {
		sRand, _ := utils.GenerateRandomString(40, "0123456789abcdefghijklmnopqrstuvwxyz");
		c.SetCookie("auth2", sRand, 2000000000, "/", "", true, true);
	}
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
