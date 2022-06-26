package auth

import (
	"sync"
	"../../settings"
)

type EntSession struct {
	SteamID64	string
	Since		int64 //unix timestamp in milliseconds
}

var MapSessions map[string]EntSession = make(map[string]EntSession);
var MuSessions sync.Mutex;


func GetSession(sSessID string) (EntSession, bool) {
	MuSessions.Lock();
	if _, ok := MapSessions[sSessID]; !ok {
		MuSessions.Unlock();
		return EntSession{}, false;
	}
	oSession := MapSessions[sSessID];
	MuSessions.Unlock();
	return oSession, true;
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
