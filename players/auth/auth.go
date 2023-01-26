package auth

import (
	"sync"
	"../../settings"
	"../../database"
	"time"
	"crypto/sha256"
	"encoding/base64"
)

type EntSession struct {
	SteamID64	string
	Since		int64 //unix timestamp in milliseconds
}

var MapSessions map[string]EntSession = make(map[string]EntSession);
var MuSessions sync.RWMutex;

func Watchers() {
	go WatchAuthExpire();
}


func WatchAuthExpire() {
	for {

		var arDeleteSessions []string;

		MuSessions.RLock();
		i64CurTime := time.Now().UnixMilli();
		for sSessID, oSession := range MapSessions {
			if (oSession.Since + settings.PlayerAuthExpire <= i64CurTime) {
				arDeleteSessions = append(arDeleteSessions, sSessID);
			}
		}
		MuSessions.RUnlock();
		for _, sSessID := range arDeleteSessions {
			RemoveSession(sSessID);
		}

		time.Sleep(86400 * time.Second); //once per day
	}
}


func GetSession(sSessID string, sCSRF string) (EntSession, bool) {
	MuSessions.RLock();
	defer MuSessions.RUnlock();

	oSession, bFound := MapSessions[sSessID];
	if (bFound) {
		by32Buffer := sha256.Sum256([]byte(sSessID));
		byBuffer := by32Buffer[:];
		if (sCSRF == base64.StdEncoding.EncodeToString(byBuffer)) {
			return oSession, true;
		}
	}
	return EntSession{}, false;
}


func GetSessionNoCSRF(sSessID string) (EntSession, bool) {
	MuSessions.RLock();
	defer MuSessions.RUnlock();

	oSession, bFound := MapSessions[sSessID];
	if (bFound) {
		return oSession, true;
	}
	return EntSession{}, false;
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
	if (sKey == "") {
		return false;
	}
	_, bKeyFound := settings.MapBackendAuthKeys[sKey];
	if (bKeyFound) {
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
