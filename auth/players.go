package auth

import (
	"sync"
	"../utils"
	"../globals"
	"time"
)

type EntSession struct {
	Player		*globals.EntPlayer
	Since		int64 //unix timestamp in milliseconds
}

var MapSessions map[string]EntSession = make(map[string]EntSession);
var MuSessions sync.Mutex;


func AddPlayerAuth(sSteamID64 string, sNicknameBase64 string) string {

	//Register player if does not exist
	globals.MuPlayers.Lock();
	if _, ok := globals.MapPlayers[sSteamID64]; !ok {

		pPlayer := &globals.EntPlayer{
			SteamID64:		sSteamID64,
		};
		globals.MapPlayers[sSteamID64] = pPlayer;
		globals.I64LastPlayerlistUpdate = time.Now().UnixMilli();

		globals.MuPlayers.Unlock();
	} else {
		globals.MuPlayers.Unlock();
	}

	sSessionKey, _ := utils.GenerateRandomString(32);

	MuSessions.Lock();

	globals.MuPlayers.Lock();
	pPlayer := globals.MapPlayers[sSteamID64];
	if (pPlayer.NicknameBase64 != sNicknameBase64) {
		pPlayer.NicknameBase64 = sNicknameBase64;
		globals.I64LastPlayerlistUpdate = time.Now().UnixMilli();
	}
	globals.MuPlayers.Unlock();

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
