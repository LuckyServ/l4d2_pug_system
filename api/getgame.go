package api

import (
	"github.com/gin-gonic/gin"
	"../settings"
	"../players/auth"
	"../players"
	"../games"
	"time"
	"fmt"
)


type GameResponse struct {
	PlayersA			[]PlayerResponse	`json:"players_a"`
	PlayersB			[]PlayerResponse	`json:"players_b"`
	GameConfig			string				`json:"game_config"`
	CampaignName		string				`json:"campaign_name"`
	PingsRequested		bool				`json:"pings_requested"`
	ServerIP			string				`json:"server_ip"`
	MmrMin				int					`json:"mmr_min"`
	MmrMax				int					`json:"mmr_max"`
	Status				string				`json:"status"`
}

var mapGameStatus = map[int]string{
	games.StateCreating:			"Creating game",
	games.StateCreated:				"Game created",
	games.StateCampaignChosen:		"Campaign selected",
	games.StateTeamsPicked:			"Players paired",
	games.StateWaitPings:			"Pinging servers",
	games.StateSelectServer:		"Selecting best available server",
	games.StateNoServers:			"No free servers available. If no server found in 5 minutes, the game ends.",
	games.StateWaitPlayersJoin:		"The server is ready. You have 5 minutes to join the server and Ready Up.",
	games.StateReadyUpExpired:		"Some players failed to Ready Up in time",
	games.StateGameProceeds:		"In game",
	games.StateGameEnded:			"Game ended",
}


func HttpReqGetGame(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.IsInGame) {

				mapResponse["success"] = true;
				games.MuGames.Lock();
				pGame := games.MapGames[pPlayer.GameID];

				var arPlayersA, arPlayersB []PlayerResponse;
				for _, pGamePlayer := range pGame.PlayersA {
					arPlayersA = append(arPlayersA, PlayerResponse{
						SteamID64:		pGamePlayer.SteamID64,
						NicknameBase64:	pGamePlayer.NicknameBase64,
						Mmr:			pGamePlayer.Mmr,
						Access:			pGamePlayer.Access,
						IsInGame:		pGamePlayer.IsInGame,
						MmrCertain:		(pGamePlayer.MmrUncertainty <= settings.MmrStable),
						IsInLobby:		pGamePlayer.IsInLobby,
						IsIdle:			pGamePlayer.IsIdle,
					});
				}
				for _, pGamePlayer := range pGame.PlayersB {
					arPlayersB = append(arPlayersB, PlayerResponse{
						SteamID64:		pGamePlayer.SteamID64,
						NicknameBase64:	pGamePlayer.NicknameBase64,
						Mmr:			pGamePlayer.Mmr,
						Access:			pGamePlayer.Access,
						IsInGame:		pGamePlayer.IsInGame,
						MmrCertain:		(pGamePlayer.MmrUncertainty <= settings.MmrStable),
						IsInLobby:		pGamePlayer.IsInLobby,
						IsIdle:			pGamePlayer.IsIdle,
					});
				}

				mapResponse["game"] = GameResponse{
					PlayersA:			arPlayersA,
					PlayersB:			arPlayersB,
					GameConfig:			pGame.GameConfig.Name,
					CampaignName:		pGame.CampaignName,
					PingsRequested:		(pGame.State == games.StateWaitPings),
					ServerIP:			pGame.ServerIP,
					MmrMin:				pGame.MmrMin,
					MmrMax:				pGame.MmrMax,
					Status:				mapGameStatus[pGame.State],
				};

				games.MuGames.Unlock();
			}
			players.MuPlayers.Unlock();
		}
	}
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.SetCookie("game_updated_at", fmt.Sprintf("%d", time.Now().UnixMilli()), 2592000, "/", "", true, false);
	c.JSON(200, mapResponse);
}
