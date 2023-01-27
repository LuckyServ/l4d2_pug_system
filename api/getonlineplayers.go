package api

import (
	"github.com/gin-gonic/gin"
	"../players"
	"../settings"
	//"fmt"
	"time"
	"../players/auth"
)


func HttpReqGetOnlinePlayers(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["authorized"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			mapResponse["authorized"] = true;
			mapResponse["steamid64"] = oSession.SteamID64;

			players.MuPlayers.RLock();

			pPlayer := players.MapPlayers[oSession.SteamID64];

			mapResponse["me"] = PlayerResponseMe{
				SteamID64:				pPlayer.SteamID64,
				NicknameBase64:			pPlayer.NicknameBase64,
				AvatarSmall:			pPlayer.AvatarSmall,
				AvatarBig:				pPlayer.AvatarBig,
				Mmr:					pPlayer.Mmr,
				Access:					pPlayer.Access,
				BanReason:				pPlayer.BanReason,
				BanAcceptedAt:			pPlayer.BanAcceptedAt,
				BanLength:				pPlayer.BanLength,
				IsInGame:				pPlayer.IsInGame,
				IsInQueue:				pPlayer.IsInQueue,
				ProfValidated:			pPlayer.ProfValidated,
				RulesAccepted:			pPlayer.RulesAccepted,
				MmrGrade:				players.GetMmrGrade(pPlayer),
				DuoOffer:				pPlayer.DuoOffer,
				IsInDuo:				(pPlayer.DuoWith != ""),
				CustomMapsState:		players.CustomMapsConfirmState(pPlayer),
			};

			players.MuPlayers.RUnlock();

		}
	}

	var arPlayers []PlayerResponse;
	var iActiveCount, iOnlineCount, iInQueueCount, iInGameCount int;

	//iStartTime := time.Now().UnixNano();
	players.MuPlayers.RLock();

	i64CurTime := time.Now().UnixMilli();
	for _, pPlayer := range players.ArrayPlayers {
		if ((pPlayer.IsOnline || pPlayer.IsInGame || pPlayer.IsInQueue) && pPlayer.ProfValidated && pPlayer.RulesAccepted && pPlayer.Access >= -1/*not banned*/) {
			arPlayers = append(arPlayers, PlayerResponse{
				SteamID64:				pPlayer.SteamID64,
				NicknameBase64:			pPlayer.NicknameBase64,
				AvatarSmall:			pPlayer.AvatarSmall,
				Mmr:					pPlayer.Mmr,
				Access:					pPlayer.Access,
				IsInGame:				pPlayer.IsInGame,
				IsInQueue:				pPlayer.IsInQueue,
				MmrGrade:				players.GetMmrGrade(pPlayer),
				IsInDuo:				(pPlayer.DuoWith != ""),
				CustomMapsState:		players.CustomMapsConfirmState(pPlayer),
			});
			if (pPlayer.IsInGame) {
				iInGameCount++;
			} else if (pPlayer.IsInQueue) {
				iInQueueCount++;
			} else if (pPlayer.IsOnline) {
				iOnlineCount++;
			}
		}
	}
	players.MuPlayers.RUnlock();
	
	iActiveCount = iOnlineCount + iInQueueCount + iInGameCount;


	mapResponse["success"] = true;
	mapResponse["count"] = map[string]int{"online": iActiveCount, "in_queue": iInQueueCount, "in_game": iInGameCount};
	mapResponse["list"] = arPlayers;
	mapResponse["newest_map"] = settings.NewestCustomMap;

	mapResponse["updated_at"] = i64CurTime;

	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	//c.SetCookie("players_updated_at", fmt.Sprintf("%d", i64CurTime), 2592000, "/", "", true, false);
	c.JSON(200, mapResponse);
}
