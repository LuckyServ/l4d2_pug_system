package api

import (
	"github.com/gin-gonic/gin"
	"../players"
	"../utils"
	"../settings"
	"fmt"
	"time"
)

type PlayerResponse struct {
	SteamID64		string		`json:"steamid64"`
	NicknameBase64	string		`json:"nickname_base64"`
	Mmr				int			`json:"mmr"`
	Access			int 		`json:"access"` //-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	IsInGame		bool		`json:"is_ingame"`
	IsIdle			bool		`json:"is_idle"`
	IsInLobby		bool		`json:"is_inlobby"`
	MmrCertain		bool		`json:"mmr_certain"`
}


func HttpReqGetOnlinePlayers(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	var arPlayers []PlayerResponse;
	var iActiveCount, iOnlineCount, iInLobbyCount, iInGameCount, iIdleCount int;

	players.MuPlayers.Lock();
	i64CurTime := time.Now().UnixMilli();
	for _, pPlayer := range players.ArrayPlayers {
		if ((pPlayer.IsOnline || pPlayer.IsInGame || pPlayer.IsInLobby) && pPlayer.ProfValidated && pPlayer.RulesAccepted && pPlayer.Access >= -1/*not banned*/) {
			bIdle := false;
			if (pPlayer.IsOnline && !pPlayer.IsInGame && !pPlayer.IsInLobby) {
				iLastAction := utils.MaxValInt64(pPlayer.OnlineSince, pPlayer.LastLobbyActivity);
				if ((i64CurTime - iLastAction) >= settings.IdleTimeout) {
					bIdle = true;
				}
			}
			arPlayers = append(arPlayers, PlayerResponse{
				SteamID64:		pPlayer.SteamID64,
				NicknameBase64:	pPlayer.NicknameBase64,
				Mmr:			pPlayer.Mmr,
				Access:			pPlayer.Access,
				IsInGame:		pPlayer.IsInGame,
				IsInLobby:		pPlayer.IsInLobby,
				IsIdle:			bIdle,
			});
			if (pPlayer.IsInGame) {
				iInGameCount++;
			} else if (pPlayer.IsInLobby) {
				iInLobbyCount++;
			} else if (pPlayer.IsOnline) {
				iOnlineCount++;
				if (bIdle) {
					iIdleCount++;
				}
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
	mapResponse["count"] = map[string]int{"online": iActiveCount, "in_lobby": iInLobbyCount, "in_game": iInGameCount, "idle": iIdleCount};
	mapResponse["list"] = arPlayers;

	
	c.Header("Access-Control-Allow-Origin", "*");
	c.SetCookie("players_updated_at", fmt.Sprintf("%d", i64CurTime), 2592000, "/", "", true, false);
	c.JSON(200, mapResponse);
}
