package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../smurf"
	"../bans"
	"../utils"
	"../database"
)


type PlayerStatusResponse struct {
	NicknameBase64			string			`json:"nickname_base64"`
	Mmr						int				`json:"mmr"`
	MmrUncertainty			float32			`json:"mmr_uncertainty"`
	LastGameResult			int				`json:"last_game_result"` //0 - unknown, 1 - draw, 2 - lost, 3 - won
	Access					int				`json:"access"`
	ProfValidated			bool			`json:"prof_validated"`
	RulesAccepted			bool			`json:"rules_accepted"`
	LastActivity			int64			`json:"last_activity"`
	IsOnline				bool			`json:"is_online"`
	OnlineSince				int64			`json:"online_since"`
	IsInGame				bool			`json:"is_ingame"`
	IsInQueue				bool			`json:"is_inqueue"`
	InQueueSince			int64			`json:"in_queue_since"`
	GameID					string			`json:"game_id"`
	GameServerPingsStored	map[string]int	`json:"gameserver_pings"`
}


func HttpReqAdminInfoOnPlayer(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	sSteamID64 := c.Query("steamid64");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			players.MuPlayers.RLock();
			iAdminAccess := players.MapPlayers[oSession.SteamID64].Access;
			players.MuPlayers.RUnlock();
			if (iAdminAccess > 0) {
				players.MuPlayers.RLock();
				pPlayer := players.MapPlayers[sSteamID64];
				if (pPlayer != nil) {

					//Current player status (all info from players.EntPlayer)
					mapResponse["status"] = PlayerStatusResponse{
						NicknameBase64:				pPlayer.NicknameBase64,
						Mmr:						pPlayer.Mmr,
						MmrUncertainty:				pPlayer.MmrUncertainty,
						LastGameResult:				pPlayer.LastGameResult,
						Access:						pPlayer.Access,
						ProfValidated:				pPlayer.ProfValidated,
						RulesAccepted:				pPlayer.RulesAccepted,
						LastActivity:				pPlayer.LastActivity,
						IsOnline:					pPlayer.IsOnline,
						OnlineSince:				pPlayer.OnlineSince,
						IsInGame:					pPlayer.IsInGame,
						IsInQueue:					pPlayer.IsInQueue,
						InQueueSince:				pPlayer.InQueueSince,
						GameID:						pPlayer.GameID,
						GameServerPingsStored:		pPlayer.GameServerPingsStored,
					};
					players.MuPlayers.RUnlock();
	
					//List of smurfs
					arSmurfs := smurf.GetKnownAccounts(sSteamID64); //slow, makes connections outside
					mapResponse["accounts"] = arSmurfs;


					//List of bans of all smurfs
					var arFilteredBanRecords []BanRecordResponse;
					bans.ChanLock <- true;

					iBanRecordsCount := len(bans.ArrayBanRecords);
					for i := iBanRecordsCount - 1; i >= 0; i-- {
						if (utils.GetStringIdxInArray(bans.ArrayBanRecords[i].SteamID64, arSmurfs) != -1) {
							arFilteredBanRecords = append(arFilteredBanRecords, BanRecordResponse{
								NicknameBase64:		bans.ArrayBanRecords[i].NicknameBase64,
								SteamID64:			bans.ArrayBanRecords[i].SteamID64,
								CreatedAt:			bans.ArrayBanRecords[i].CreatedAt,
								BannedBySteamID64:	bans.ArrayBanRecords[i].BannedBySteamID64,
								AcceptedAt:			bans.ArrayBanRecords[i].AcceptedAt,
								BanLength:			bans.ArrayBanRecords[i].BanLength,
								BanReasonBase64:	bans.ArrayBanRecords[i].BanReasonBase64,
							});
							if (len(arFilteredBanRecords) > 500) {
								break;
							}
						}
					}

					bans.ChanUnlock <- true;

					mapResponse["bans"] = arFilteredBanRecords;


					//List of used IP addresses and locations (hide IP addresses if access < 3)
					mapResponse["ips"] = smurf.GetIPInfo(sSteamID64);
					

					//Game history
					arGameHistory, iGamesPlayed := database.GetGameHistory(sSteamID64);
					mapGameHistory := make(map[string]interface{});
					mapGameHistory["count"] = iGamesPlayed;
					mapGameHistory["history"] = arGameHistory;
					mapResponse["games"] = mapGameHistory;

					mapResponse["success"] = true;
				} else {
					players.MuPlayers.RUnlock();
					mapResponse["error"] = "Player not found";
				}
			} else {
				mapResponse["error"] = "You dont have access to this information";
			}
		} else {
			mapResponse["error"] = "Please authorize first";
		}
	} else {
		mapResponse["error"] = "Please authorize first";
	}
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
