package api

import (
	"github.com/gin-gonic/gin"
	"../players"
	"fmt"
	"time"
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
