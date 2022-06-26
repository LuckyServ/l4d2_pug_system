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
	Pings			map[string]int
	LastPingsUpdate	int64 //unix timestamp in milliseconds
	LastActivity	int64 //unix timestamp in milliseconds
	IsOnline		bool
	IsInGame		bool
	IsInLobby		bool
	LastChanged		int64 //Last time player info was changed //unix timestamp in milliseconds
}

var MapPlayers map[string]*EntPlayer = make(map[string]*EntPlayer);
var MuPlayers sync.Mutex;

var I64LastPlayerlistUpdate int64;


func UpdatePlayerActivity(sSteamID64 string) {
	MuPlayers.Lock();
	if _, ok := MapPlayers[sSteamID64]; !ok {
		MuPlayers.Unlock();
		return;
	}
	MapPlayers[sSteamID64].LastActivity = time.Now().UnixMilli();
	MuPlayers.Unlock();
}


func AddPlayerAuth(sSteamID64 string, sNicknameBase64 string) string {

	//Register player if does not exist
	MuPlayers.Lock();
	if _, ok := MapPlayers[sSteamID64]; !ok {

		pPlayer := &EntPlayer{
			SteamID64:			sSteamID64,
			MmrUncertainty:		settings.DefaultMmrUncertainty,
		};
		MapPlayers[sSteamID64] = pPlayer;
		I64LastPlayerlistUpdate = time.Now().UnixMilli();

		MuPlayers.Unlock();
	} else {
		MuPlayers.Unlock();
	}

	sSessionKey, _ := utils.GenerateRandomString(32);

	auth.MuSessions.Lock();

	MuPlayers.Lock();
	if (MapPlayers[sSteamID64].NicknameBase64 != sNicknameBase64) {
		MapPlayers[sSteamID64].NicknameBase64 = sNicknameBase64;
		I64LastPlayerlistUpdate = time.Now().UnixMilli();
	}
	MuPlayers.Unlock();

	oSession := auth.EntSession{
		SteamID64:	sSteamID64,
		Since:		time.Now().UnixMilli(),
	};

	auth.MapSessions[sSessionKey] = oSession;
	auth.MuSessions.Unlock();

	return sSessionKey;
}
