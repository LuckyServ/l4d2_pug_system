package auth

import (
	"sync"
	"../utils"
	"../players"
	"time"
)

type EntSession struct {
	Player		*players.EntPlayer
	Since		int64 //unix timestamp in milliseconds
}

var MapSessions map[string]EntSession = make(map[string]EntSession);
var MuSessions sync.Mutex;


func AddPlayerAuth(sSteamID64 string, sNicknameBase64 string) string {

	//Register player if does not exist
	players.MuPlayers.Lock();
	if _, ok := players.MapPlayers[sSteamID64]; !ok {

		pPlayer := &players.EntPlayer{
			SteamID64:		sSteamID64,
		};
		players.MapPlayers[sSteamID64] = pPlayer;
		players.I64LastPlayerlistUpdate = time.Now().UnixMilli();

		players.MuPlayers.Unlock();
	} else {
		players.MuPlayers.Unlock();
	}

	sSessionKey, _ := utils.GenerateRandomString(32);

	MuSessions.Lock();

	players.MuPlayers.Lock();
	pPlayer := players.MapPlayers[sSteamID64];
	if (pPlayer.NicknameBase64 != sNicknameBase64) {
		pPlayer.NicknameBase64 = sNicknameBase64;
		players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
	}
	players.MuPlayers.Unlock();

	oSession := EntSession{
		Player:		pPlayer,
		Since:		time.Now().UnixMilli(),
	};

	MapSessions[sSessionKey] = oSession;
	MuSessions.Unlock();

	return sSessionKey;
}

func RemovePlayerAuth(sSessID string) (bool, string) {

	MuSessions.Lock();
	if _, ok := MapSessions[sSessID]; !ok {
		MuSessions.Unlock();
		return false, "Session ID does not exist";
	}
	delete(MapSessions, sSessID);
	MuSessions.Unlock();

	return true, "";
}
