package players

import (
	"sync"
	"time"
	"../settings"
	"../utils"
	"./auth"
)

type EntPlayer struct {
	SteamID64		string
	NicknameBase64	string
	Mmr				int
	MmrUncertainty	int
	Access			int //-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	ProfValidated	bool //Steam profile validated
	RulesAccepted	bool //Rules accepted
	LastActivity	int64 //unix timestamp in milliseconds
	IsOnline		bool
	IsInGame		bool
	IsInLobby		bool
	LastChanged		int64 //Last time player info was changed //unix timestamp in milliseconds
}

var MapPlayers map[string]*EntPlayer = make(map[string]*EntPlayer);
var MuPlayers sync.Mutex;

var I64LastPlayerlistUpdate int64;


func UpdatePlayerActivity(sSteamID64 string) { //Maps must be locked outside!!!
	if _, ok := MapPlayers[sSteamID64]; !ok {
		return;
	}
	i64CurTime := time.Now().UnixMilli();
	MapPlayers[sSteamID64].LastActivity = i64CurTime;
	if (MapPlayers[sSteamID64].IsOnline == false) {
		MapPlayers[sSteamID64].IsOnline = true;
		MapPlayers[sSteamID64].LastChanged = i64CurTime;
		I64LastPlayerlistUpdate = i64CurTime;
	}
}


func AddPlayerAuth(sSteamID64 string, sNicknameBase64 string) string {
	i64CurTime := time.Now().UnixMilli();

	//Register player if does not exist
	MuPlayers.Lock();
	if _, ok := MapPlayers[sSteamID64]; !ok {

		pPlayer := &EntPlayer{
			SteamID64:			sSteamID64,
			MmrUncertainty:		settings.DefaultMmrUncertainty,
			LastChanged:		i64CurTime,
		};
		MapPlayers[sSteamID64] = pPlayer;
		I64LastPlayerlistUpdate = i64CurTime;

		MuPlayers.Unlock();
	} else {
		MuPlayers.Unlock();
	}

	sSessionKey, _ := utils.GenerateRandomString(32);

	auth.MuSessions.Lock();

	MuPlayers.Lock();
	if (MapPlayers[sSteamID64].NicknameBase64 != sNicknameBase64) {
		MapPlayers[sSteamID64].NicknameBase64 = sNicknameBase64;
		I64LastPlayerlistUpdate = i64CurTime;
	}
	MuPlayers.Unlock();

	oSession := auth.EntSession{
		SteamID64:	sSteamID64,
		Since:		i64CurTime,
	};

	auth.MapSessions[sSessionKey] = oSession;
	auth.MuSessions.Unlock();

	return sSessionKey;
}
