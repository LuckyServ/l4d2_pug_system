package players

import (
	"sync"
	"time"
	"../settings"
)

type EntPlayer struct {
	SteamID64		string
	NicknameBase64	string
	Mmr				int
	MmrUncertainty	int
	Access			int //-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	ProfValidated	bool //Steam profile validated
	Pings			map[string]int
	LastPingsUpdate	int64 //unix timestamp in milliseconds
	LastActivity	int64 //unix timestamp in milliseconds
	IsOnline		bool
	IsInGame		bool
	IsInLobby		bool
	LastUpdated		int64 //Last time player info was changed //unix timestamp in milliseconds
}

type SinglePlayerHTTPResponse struct {
	NicknameBase64	string
	Mmr				int
	MmrUncertainty	int
	MmrCertain		bool
	Access			int //-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	ProfValidated	bool //Steam profile validated
	PingsUpdated	bool //Do we have all servers pinged
	IsOnline		bool
	IsInGame		bool
	IsInLobby		bool
	LastUpdated		int64 //Last time player info was changed //unix timestamp in milliseconds
}

var MapPlayers map[string]*EntPlayer = make(map[string]*EntPlayer);
var MuPlayers sync.Mutex;

var I64LastPlayerlistUpdate int64;


func UpdatePlayerActivity(sSteamID64 string) (bool, int) {
	MuPlayers.Lock();
	if _, ok := MapPlayers[sSteamID64]; !ok {
		MuPlayers.Unlock();
		return false, 3;
	}
	MapPlayers[sSteamID64].LastActivity = time.Now().UnixMilli();
	MuPlayers.Unlock();
	return true, 0;
}

func GetPlayerResponse(sSteamID64 string) (bool, SinglePlayerHTTPResponse, int) {

	MuPlayers.Lock();
	if _, ok := MapPlayers[sSteamID64]; !ok {
		MuPlayers.Unlock();
		return false, SinglePlayerHTTPResponse{}, 3;
	}

	var bMmrCertain, bPingsUpdated bool;
	if (MapPlayers[sSteamID64].MmrUncertainty <= settings.MmrStable) {
		bMmrCertain = true;
	} else {
		bMmrCertain = false;
	}
	if ((time.Now().UnixMilli() - MapPlayers[sSteamID64].LastPingsUpdate) < settings.PingsMaxAge) {
		bPingsUpdated = true;
	} else {
		bPingsUpdated = false;
	}

	oPlayer := SinglePlayerHTTPResponse{
		NicknameBase64:	MapPlayers[sSteamID64].NicknameBase64,
		Mmr:			MapPlayers[sSteamID64].Mmr,
		MmrUncertainty:	MapPlayers[sSteamID64].MmrUncertainty,
		MmrCertain:		bMmrCertain,
		Access:			MapPlayers[sSteamID64].Access,
		ProfValidated:	MapPlayers[sSteamID64].ProfValidated,
		PingsUpdated:	bPingsUpdated,
		IsOnline:		MapPlayers[sSteamID64].IsOnline,
		IsInGame:		MapPlayers[sSteamID64].IsInGame,
		IsInLobby:		MapPlayers[sSteamID64].IsInLobby,
		LastUpdated:	MapPlayers[sSteamID64].LastUpdated,
	};

	MuPlayers.Unlock();
	return true, oPlayer, 0;
}
