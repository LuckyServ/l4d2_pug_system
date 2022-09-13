package smurf

import (
	"time"
	"sync"
)

type EntVPNInfo struct {
	IsVPN		bool
	UpdatedAt	int64 //unix time in seconds
}

var mapVPNs = make(map[string]EntVPNInfo);
var MuVPN sync.RWMutex;


func Watchers() {
	for {
		time.Sleep(86400 * time.Second); //24 hours
		var arRemoveIP []string;
		MuVPN.RLock();
		i64CurTime := time.Now().Unix();
		for sIP, oVPNInfo := range mapVPNs {
			if (oVPNInfo.UpdatedAt + 604800 <= i64CurTime) {
				arRemoveIP = append(arRemoveIP, sIP);
			}
		}
		MuVPN.RUnlock();
		MuVPN.Lock();
		for _, sIP := range arRemoveIP {
			delete(mapVPNs, sIP);
		}
		MuVPN.Unlock();
	}
}

func AnnounceIP(sIP string) { //thread safe, fast
}

func IsVPN(sIP string) bool { //thread safe, fast
	return false;
}
