package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../games"
	"../settings"
	"time"
	//"fmt"
)


type GameResponse struct {
	ID					string				`json:"id"`
	PlayersA			[]PlayerResponse	`json:"players_a"`
	PlayersB			[]PlayerResponse	`json:"players_b"`
	GameConfig			string				`json:"game_config"`
	CampaignName		string				`json:"campaign_name"`
	PingsRequested		bool				`json:"pings_requested"`
	ServerIP			string				`json:"server_ip"`
	ProxyIP				string				`json:"proxy_ip"`
	Status				string				`json:"status"`
	MapDownloadLink		string				`json:"map_download"`
}


func HttpReqGetGame(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			players.MuPlayers.RLock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.IsInGame) {

				mapResponse["success"] = true;
				games.MuGames.RLock();
				pGame := games.MapGames[pPlayer.GameID];

				var arPlayersA, arPlayersB []PlayerResponse;
				for _, pGamePlayer := range pGame.PlayersA {
					arPlayersA = append(arPlayersA, PlayerResponse{
						SteamID64:			pGamePlayer.SteamID64,
						NicknameBase64:		pGamePlayer.NicknameBase64,
						AvatarSmall:		pGamePlayer.AvatarSmall,
						Mmr:				pGamePlayer.Mmr,
						Access:				pGamePlayer.Access,
						IsInGame:			pGamePlayer.IsInGame,
						MmrGrade:			players.GetMmrGrade(pGamePlayer),
						IsInQueue:			pGamePlayer.IsInQueue,
						IsInDuo:			(pGamePlayer.DuoWith != ""),
						CustomMapsState:	players.CustomMapsConfirmState(pGamePlayer),
					});
				}
				for _, pGamePlayer := range pGame.PlayersB {
					arPlayersB = append(arPlayersB, PlayerResponse{
						SteamID64:			pGamePlayer.SteamID64,
						NicknameBase64:		pGamePlayer.NicknameBase64,
						AvatarSmall:		pGamePlayer.AvatarSmall,
						Mmr:				pGamePlayer.Mmr,
						Access:				pGamePlayer.Access,
						IsInGame:			pGamePlayer.IsInGame,
						MmrGrade:			players.GetMmrGrade(pGamePlayer),
						IsInQueue:			pGamePlayer.IsInQueue,
						IsInDuo:			(pGamePlayer.DuoWith != ""),
						CustomMapsState:	players.CustomMapsConfirmState(pGamePlayer),
					});
				}

				mapResponse["game"] = GameResponse{
					ID:					pGame.ID,
					PlayersA:			arPlayersA,
					PlayersB:			arPlayersB,
					GameConfig:			pGame.GameConfig.Name,
					CampaignName:		pGame.CampaignName,
					PingsRequested:		(pGame.State == games.StateWaitPings),
					ServerIP:			pGame.ServerIP,
					ProxyIP:			settings.MapProxies[pGame.ServerIP],
					Status:				games.MapGameStatus[pGame.State],
					MapDownloadLink:	pGame.MapDownloadLink,
				};

				games.MuGames.RUnlock();
			}
			players.MuPlayers.RUnlock();
		}
	}

	mapResponse["updated_at"] = time.Now().UnixMilli();
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	//c.SetCookie("game_updated_at", fmt.Sprintf("%d", time.Now().UnixMilli()), 2592000, "/", "", true, false);
	c.JSON(200, mapResponse);
}
