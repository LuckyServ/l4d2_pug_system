package api

import (
	"github.com/gin-gonic/gin"
	"../lobby"
	"fmt"
	"time"
)

type LobbyResponse struct {
	ID				string		`json:"id"`
	MmrMin			int			`json:"mmr_min"`
	MmrMax			int			`json:"mmr_max"`
	CreatedAt		int64		`json:"created_at"` //milliseconds
	GameConfig		string		`json:"confogl_config"`
	PlayerCount		int			`json:"player_count"`
	ReadyUpState	bool		`json:"readyup_state"`
}


func HttpReqGetLobbies(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	var arRespLobbies []LobbyResponse;
	var iLobbiesCount int;

	lobby.MuLobbies.Lock();

	for _, pLobby := range lobby.ArrayLobbies {
		arRespLobbies = append(arRespLobbies, LobbyResponse{
			ID:				pLobby.ID,
			MmrMin:			pLobby.MmrMin,
			MmrMax:			pLobby.MmrMax,
			CreatedAt:		pLobby.CreatedAt,
			GameConfig:		pLobby.GameConfig,
			PlayerCount:	pLobby.PlayerCount,
			ReadyUpState:	pLobby.ReadyUpState,
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
	mapResponse["lobbies"] = arRespLobbies;

	
	c.Header("Access-Control-Allow-Origin", "*");
	c.SetCookie("lobbies_updated_at", fmt.Sprintf("%d", time.Now().UnixMilli()), 2592000, "/", "", true, false);
	c.JSON(200, mapResponse);
}
