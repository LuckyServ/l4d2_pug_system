package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../queue"
	"../games"
	//"fmt"
	"time"
)


func HttpReqGetQueue(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["authorized"] = false;
	mapResponse["is_inqueue"] = false;
	mapResponse["need_readyup"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			mapResponse["authorized"] = true;
			players.MuPlayers.RLock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.IsInQueue) {
				mapResponse["is_inqueue"] = true;
				if (pPlayer.IsReadyUpRequested && !pPlayer.IsReadyConfirmed) {
					mapResponse["need_readyup"] = true;
				}
				mapResponse["duo_with_nickbase64"] = "";
				if (pPlayer.DuoWith != "") {
					pPlayer2, bFound := players.MapPlayers[pPlayer.DuoWith];
					if (bFound && pPlayer2 != nil) {
						mapResponse["duo_with_nickbase64"] = pPlayer2.NicknameBase64;
					}
				}
			}
			players.MuPlayers.RUnlock();
		}
	}


	mapResponse["success"] = true;
	mapResponse["player_count"] = queue.IPlayersCount;
	players.MuPlayers.RLock();
	if (queue.PLongestWaitPlayer != nil) {
		mapResponse["waiting_since"] = queue.PLongestWaitPlayer.InQueueSince;
		mapResponse["waiting_till"] = queue.PLongestWaitPlayer.InQueueSince + queue.I64MaxQueueWait;
	} else {
		mapResponse["waiting_since"] = 0;
		mapResponse["waiting_till"] = 0;
	}
	mapResponse["ready_players"] = queue.IReadyPlayers;
	mapResponse["ready_state"] = queue.BIsInReadyUp;
	mapResponse["finishing_game"] = games.IPlayersFinishingGameSoon;
	players.MuPlayers.RUnlock();

	mapResponse["updated_at"] = time.Now().UnixMilli();

	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	//c.SetCookie("queue_updated_at", fmt.Sprintf("%d", time.Now().UnixMilli()), 2592000, "/", "", true, false);
	c.JSON(200, mapResponse);
}
