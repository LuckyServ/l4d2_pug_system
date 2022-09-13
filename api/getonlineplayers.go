package api

import (
	"github.com/gin-gonic/gin"
	"../players"
	"fmt"
	"time"
	"../settings"
	"../players/auth"
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

type PlayerResponseMe struct {
	SteamID64		string		`json:"steamid64"`
	NicknameBase64	string		`json:"nickname_base64"`
	Mmr				int			`json:"mmr"`
	Access			int 		`json:"access"` //-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	IsInGame		bool		`json:"is_ingame"`
	IsIdle			bool		`json:"is_idle"`
	IsInLobby		bool		`json:"is_inlobby"`
	MmrCertain		bool		`json:"mmr_certain"`
	ProfValidated	bool		`json:"profile_validated"` //Steam profile validated
	RulesAccepted	bool		`json:"rules_accepted"` //Rules accepted
}


func HttpReqGetOnlinePlayers(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["authorized"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			mapResponse["authorized"] = true;
			mapResponse["steamid64"] = oSession.SteamID64;

			players.MuPlayers.RLock();

			pPlayer := players.MapPlayers[oSession.SteamID64];

			mapResponse["me"] = PlayerResponseMe{
				SteamID64:		pPlayer.SteamID64,
				NicknameBase64:	pPlayer.NicknameBase64,
				Mmr:			pPlayer.Mmr,
				Access:			pPlayer.Access,
				IsInGame:		pPlayer.IsInGame,
				IsInLobby:		pPlayer.IsInLobby,
				IsIdle:			pPlayer.IsIdle,
				MmrCertain:		(pPlayer.MmrUncertainty <= settings.MmrStable),
				ProfValidated:	pPlayer.ProfValidated,
				RulesAccepted:	pPlayer.RulesAccepted,
			};

			players.MuPlayers.RUnlock();

		}
	}

	var arPlayers []PlayerResponse;
	var iActiveCount, iOnlineCount, iInLobbyCount, iInGameCount, iIdleCount int;

	//iStartTime := time.Now().UnixNano();
	players.MuPlayers.RLock();

	i64CurTime := time.Now().UnixMilli();
	for _, pPlayer := range players.ArrayPlayers {
		if ((pPlayer.IsOnline || pPlayer.IsInGame || pPlayer.IsInLobby) && pPlayer.ProfValidated && pPlayer.RulesAccepted && pPlayer.Access >= -1/*not banned*/) {
			arPlayers = append(arPlayers, PlayerResponse{
				SteamID64:		pPlayer.SteamID64,
				NicknameBase64:	pPlayer.NicknameBase64,
				Mmr:			pPlayer.Mmr,
				Access:			pPlayer.Access,
				IsInGame:		pPlayer.IsInGame,
				MmrCertain:		(pPlayer.MmrUncertainty <= settings.MmrStable),
				IsInLobby:		pPlayer.IsInLobby,
				IsIdle:			pPlayer.IsIdle,
			});
			if (pPlayer.IsInGame) {
				iInGameCount++;
			} else if (pPlayer.IsInLobby) {
				iInLobbyCount++;
			} else if (pPlayer.IsOnline) {
				iOnlineCount++;
				if (pPlayer.IsIdle) {
					iIdleCount++;
				}
			}
		}
	}
	players.MuPlayers.RUnlock();
	//fmt.Printf("Was locked for %d μs\n", (time.Now().UnixNano() - iStartTime) / 1000);
	iActiveCount = iOnlineCount + iInLobbyCount + iInGameCount;


	mapResponse["success"] = true;
	mapResponse["count"] = map[string]int{"online": iActiveCount, "in_lobby": iInLobbyCount, "in_game": iInGameCount, "idle": iIdleCount};
	mapResponse["list"] = arPlayers;

	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.SetCookie("players_updated_at", fmt.Sprintf("%d", i64CurTime), 2592000, "/", "", true, false);
	c.JSON(200, mapResponse);
}
