package api

import (
	"github.com/gin-gonic/gin"
	"../players"
	"../settings"
	//"fmt"
	"time"
	"../players/auth"
)

type PlayerResponse struct {
	SteamID64				string		`json:"steamid64"`
	NicknameBase64			string		`json:"nickname_base64"`
	AvatarSmall				string		`json:"avatar_small"`
	Mmr						int			`json:"mmr"`
	Access					int 		`json:"access"` //-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	IsInGame				bool		`json:"is_ingame"`
	IsInQueue				bool		`json:"is_inqueue"`
	MmrGrade				int			`json:"mmr_grade"`
	CustomMapsState			int			`json:"custom_maps"` //1 - never confirmed, 2 - update required, 3 - confirmed
	IsInDuo					bool		`json:"in_duo"`
}

type PlayerResponseMe struct {
	SteamID64				string		`json:"steamid64"`
	NicknameBase64			string		`json:"nickname_base64"`
	AvatarSmall				string		`json:"avatar_small"`
	AvatarBig				string		`json:"avatar_big"`
	Mmr						int			`json:"mmr"`
	Access					int 		`json:"access"` //-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	BanReason				string 		`json:"banreason"`
	BanAcceptedAt			int64 		`json:"ban_accepted_at"`
	BanLength				int64 		`json:"ban_length"`
	IsInGame				bool		`json:"is_ingame"`
	IsInQueue				bool		`json:"is_inqueue"`
	MmrGrade				int			`json:"mmr_grade"`
	ProfValidated			bool		`json:"profile_validated"` //Steam profile validated
	RulesAccepted			bool		`json:"rules_accepted"` //Rules accepted
	DuoOffer				string		`json:"duo_offer"`
	CustomMapsState			int			`json:"custom_maps"` //1 - never confirmed, 2 - update required, 3 - confirmed
	IsInDuo					bool		`json:"in_duo"`
}


func HttpReqGetPlayers(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	sSearchSteamID := c.Query("steamid64");
	bOnlyOnline := c.Query("online") == "true";

	mapResponse["authorized"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			mapResponse["authorized"] = true;

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
	bAuthorized := mapResponse["authorized"] == true;

	if (bOnlyOnline) {
		for _, pPlayer := range players.ArrayPlayers {
			if (pPlayer.ProfValidated && pPlayer.RulesAccepted && pPlayer.Access >= -1 && (pPlayer.IsOnline || pPlayer.IsInGame || pPlayer.IsInQueue)) {
				arPlayers = AppendPlayerResponse(arPlayers, pPlayer);
				if (!bAuthorized && len(arPlayers) >= 15) {
					break;
				}
			}
		}
	} else if (sSearchSteamID != "") {
		pSearchedPlayer, bFound := players.MapPlayers[sSearchSteamID];
		if (bFound) {
			arPlayers = AppendPlayerResponse(arPlayers, pSearchedPlayer);
		}
	}/* else {
		for _, pPlayer := range players.ArrayPlayers {
			if (pPlayer.ProfValidated && pPlayer.RulesAccepted && pPlayer.Access >= -1 && players.GetMmrGrade(pPlayer) > 0) {
				arPlayers = AppendPlayerResponse(arPlayers, pPlayer);
				if (!bAuthorized && count(arPlayers) >= 15) {
					break;
				}
			}
		}
	}*/

	for _, pPlayer := range players.ArrayPlayers {
		if (pPlayer.ProfValidated && pPlayer.RulesAccepted && pPlayer.Access >= -1) {
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

func AppendPlayerResponse(arPlayers []PlayerResponse, pPlayer *players.EntPlayer) []PlayerResponse {
	arBuffer := append(arPlayers, PlayerResponse{
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
	return arBuffer;
}