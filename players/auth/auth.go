package auth

import (
	"sync"
	"../../settings"
	"../../database"
)

type EntSession struct {
	SteamID64	string
	Since		int64 //unix timestamp in milliseconds
}

var MapSessions map[string]EntSession = make(map[string]EntSession);
var MuSessions sync.RWMutex;


func GetSession(sSessID string) (EntSession, bool) {
	MuSessions.RLock();
	if _, ok := MapSessions[sSessID]; !ok {
		MuSessions.RUnlock();
		return EntSession{}, false;
	}
	oSession := MapSessions[sSessID];
	MuSessions.RUnlock();
	return oSession, true;
}

func RemoveSession(sSessID string) bool {
	MuSessions.Lock();
	if _, ok := MapSessions[sSessID]; !ok {
		MuSessions.Unlock();
		return false;
	}
	delete(MapSessions, sSessID);
	MuSessions.Unlock();
	go database.RemoveSession(sSessID);
	return true;
}

func RestoreSessions() bool { //no need to lock maps
	arDatabaseSessions := database.RestoreSessions();
	for _, oDBSession := range arDatabaseSessions {
		oSession := EntSession{
			SteamID64:	oDBSession.SteamID64,
			Since:		oDBSession.Since,
		};
		MapSessions[oDBSession.SessionID] = oSession;
	}
	return true;
}

func Backend(sKey string) bool {
	if (sKey == settings.BackendAuthKey) {
		return true;
	}
	return false;
}

/*func RemovePlayerAuth(sSessID string) (bool, int) {

	MuSessions.Lock();
	if _, ok := MapSessions[sSessID]; !ok {
		MuSessions.Unlock();
		return false, 3;
	}
	delete(MapSessions, sSessID);
	MuSessions.Unlock();

	return true, 0;
}*/
