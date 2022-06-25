package players

import (
	"sync"
	"time"
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
