package players

import (
	"sync"
	"time"
	"../settings"
	"../utils"
	"../database"
	"./auth"
	"strconv"
)

type EntPlayer struct {
	SteamID64		string
	NicknameBase64	string
	Mmr				int
	MmrUncertainty	float32
	Access			int //-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	ProfValidated	bool //Steam profile validated
	RulesAccepted	bool //Rules accepted
	LastActivity	int64 //unix timestamp in milliseconds
	IsOnline		bool
	IsInGame		bool
	IsInLobby		bool
	ReadyInLobby	bool
	LastChanged		int64 //Last time player info was changed //unix timestamp in milliseconds
	LastValidateReq	int64 //Last profile validation request //unix timestamp in milliseconds
}

var MapPlayers map[string]*EntPlayer = make(map[string]*EntPlayer);
var ArrayPlayers []*EntPlayer; //duplicate of MapPlayers, for faster iterating
var MuPlayers sync.Mutex;

var I64LastPlayerlistUpdate int64;


func UpdatePlayerActivity(sSteamID64 string) { //Maps must be locked outside!!!
	if _, ok := MapPlayers[sSteamID64]; !ok {
		return;
	}
	i64CurTime := time.Now().UnixMilli();
	pPlayer := MapPlayers[sSteamID64];
	pPlayer.LastActivity = i64CurTime;
	if (!pPlayer.IsOnline) {
		pPlayer.IsOnline = true;
		pPlayer.LastChanged = i64CurTime;
		I64LastPlayerlistUpdate = i64CurTime;
	}
}

func RestorePlayers() bool { //no need to lock maps
	arDatabasePlayers := database.RestorePlayers();
	/*if (len(arDatabasePlayers) == 0) {
		return false;
	}*/
	for _, oDBPlayer := range arDatabasePlayers {
		pPlayer := &EntPlayer{
			SteamID64:			oDBPlayer.SteamID64,
			NicknameBase64:		oDBPlayer.NicknameBase64,
			Mmr:				oDBPlayer.Mmr,
			MmrUncertainty:		oDBPlayer.MmrUncertainty,
			Access:				oDBPlayer.Access,
			ProfValidated:		oDBPlayer.ProfValidated,
			RulesAccepted:		oDBPlayer.RulesAccepted,
			LastChanged:		time.Now().UnixMilli(),
		};
		MapPlayers[oDBPlayer.SteamID64] = pPlayer;
		ArrayPlayers = append(ArrayPlayers, pPlayer);
	}
	return true;
}


func AddPlayerAuth(sSteamID64 string, sNicknameBase64 string) string {

	//Register player if does not exist
	MuPlayers.Lock();
	if _, ok := MapPlayers[sSteamID64]; !ok {

		pPlayer := &EntPlayer{
			SteamID64:			sSteamID64,
			MmrUncertainty:		settings.DefaultMmrUncertainty,
			LastChanged:		time.Now().UnixMilli(),
		};
		MapPlayers[sSteamID64] = pPlayer;
		ArrayPlayers = append(ArrayPlayers, pPlayer);
		I64LastPlayerlistUpdate = time.Now().UnixMilli();

		oDatabasePlayer := database.DatabasePlayer{
			SteamID64:			sSteamID64,
			MmrUncertainty:		settings.DefaultMmrUncertainty,
		};
		MuPlayers.Unlock();

		database.AddPlayer(oDatabasePlayer);

	} else {
		MuPlayers.Unlock();
	}

	sSessionKey, _ := utils.GenerateRandomString(32, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz");
	sSessionKey = sSessionKey+strconv.FormatInt(time.Now().UnixNano(), 10);

	auth.MuSessions.Lock();

	MuPlayers.Lock();
	pPlayer := MapPlayers[sSteamID64];
	if (pPlayer.NicknameBase64 != sNicknameBase64) {

		pPlayer.NicknameBase64 = sNicknameBase64;
		pPlayer.LastChanged = time.Now().UnixMilli();
		I64LastPlayerlistUpdate = time.Now().UnixMilli();

		go database.UpdatePlayer(database.DatabasePlayer{
			SteamID64:			pPlayer.SteamID64,
			NicknameBase64:		pPlayer.NicknameBase64,
			Mmr:				pPlayer.Mmr,
			MmrUncertainty:		pPlayer.MmrUncertainty,
			Access:				pPlayer.Access,
			ProfValidated:		pPlayer.ProfValidated,
			RulesAccepted:		pPlayer.RulesAccepted,
			});
	}
	MuPlayers.Unlock();

	oSession := auth.EntSession{
		SteamID64:	sSteamID64,
		Since:		time.Now().UnixMilli(),
	};

	auth.MapSessions[sSessionKey] = oSession;
	go database.AddSession(database.DatabaseSession{
		SessionID:			sSessionKey,
		SteamID64:			oSession.SteamID64,
		Since:				oSession.Since,
		});
	auth.MuSessions.Unlock();

	return sSessionKey;
}
