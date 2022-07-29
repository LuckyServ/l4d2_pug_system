package api

import (
	"github.com/gin-gonic/gin"
	"../lobby"
	"../players/auth"
	"../players"
	"fmt"
	"time"
	"../settings"
)

type LobbyResponse struct {
	ID				string		`json:"id"`
	MmrMin			int			`json:"mmr_min"`
	MmrMax			int			`json:"mmr_max"`
	CreatedAt		int64		`json:"created_at"` //milliseconds
	GameConfig		string		`json:"confogl_config"`
	PlayerCount		int			`json:"player_count"`
	ReadyUpState	bool		`json:"readyup_state"`
	ReadyPlayers	int			`json:"ready_players"`
}


func HttpReqGetLobbies(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	var arRespLobbies []LobbyResponse;
	var iLobbiesCount int;
	var bNeedReadyUp bool;

	lobby.MuLobbies.Lock();

	mapResponse["authorized"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			mapResponse["authorized"] = true;
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.IsInLobby) {
				pLobby := lobby.MapLobbies[pPlayer.LobbyID];
				mapResponse["mylobby"] = LobbyResponse{
					ID:				pLobby.ID,
					MmrMin:			pLobby.MmrMin,
					MmrMax:			pLobby.MmrMax,
					CreatedAt:		pLobby.CreatedAt,
					GameConfig:		pLobby.GameConfig,
					PlayerCount:	pLobby.PlayerCount,
					ReadyUpState:	(pLobby.PlayerCount >= 8),
					ReadyPlayers:	pLobby.ReadyPlayers,
				};
				if (pLobby.PlayerCount >= 8 && !pPlayer.IsReadyInLobby) {
					bNeedReadyUp = true;
				}
			}
			players.MuPlayers.Unlock();
		}
	}

	for _, pLobby := range lobby.ArrayLobbies {
		arRespLobbies = append(arRespLobbies, LobbyResponse{
			ID:				pLobby.ID,
			MmrMin:			pLobby.MmrMin,
			MmrMax:			pLobby.MmrMax,
			CreatedAt:		pLobby.CreatedAt,
			GameConfig:		pLobby.GameConfig,
			PlayerCount:	pLobby.PlayerCount,
			ReadyUpState:	(pLobby.PlayerCount >= 8),
			ReadyPlayers:	pLobby.ReadyPlayers,
		});
		iLobbiesCount++;
	}

	lobby.MuLobbies.Unlock();

	
	//sort
	iSize := len(arRespLobbies);
	if (iSize > 1) {
		bSorted := false;
		for !bSorted {
			bSorted = true;
			for i := 1; i < iSize; i++ {
				if (arRespLobbies[i].CreatedAt < arRespLobbies[i - 1].CreatedAt) {
					arRespLobbies[i], arRespLobbies[i - 1] = arRespLobbies[i - 1], arRespLobbies[i]; //switch
					bSorted = false;
				}
			}
		}
	}


	mapResponse["success"] = true;
	mapResponse["count"] = iLobbiesCount;
	mapResponse["need_readyup"] = bNeedReadyUp;
	mapResponse["lobbies"] = arRespLobbies;

	
	c.Header("Access-Control-Allow-Origin", "https://"+settings.HomeDomain);
	c.Header("Access-Control-Allow-Credentials", "true");
	c.SetCookie("lobbies_updated_at", fmt.Sprintf("%d", time.Now().UnixMilli()), 2592000, "/", "", true, false);
	c.JSON(200, mapResponse);
}
