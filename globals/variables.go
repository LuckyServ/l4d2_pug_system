package globals

import (
	"sync"
)

type EntPlayer struct {
	SteamID64		string
	NicknameBase64	string
	Mmr				int
	MmrUncertainty	int
	Access			int //-1 - banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	ProfValidated	bool //Steam profile validated
	Pings			map[string]int
	PingsUpdated	int64 //unix timestamp in milliseconds
}

var MapPlayers map[string]*EntPlayer = make(map[string]*EntPlayer);
var MuPlayers sync.Mutex;
