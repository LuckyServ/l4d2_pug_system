package smurf

import (
	"time"
	"sync"
	"io"
	"github.com/buger/jsonparser"
	"net/http"
	"fmt"
	"../settings"
)

type EntVPNInfo struct {
	IsVPN			bool
	IsInCheck		bool
	UpdatedAt		int64 //unix time in seconds
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
	MuVPN.RLock();
	bShouldCheck := true;
	oVPNInfo, bFound := mapVPNs[sIP];
	if (bFound) {
		if (oVPNInfo.IsInCheck && oVPNInfo.UpdatedAt + 60/*1min*/ > time.Now().Unix()) {
			bShouldCheck = false;
		}
	}
	MuVPN.RUnlock();

	if (bShouldCheck) {
		MuVPN.Lock();
		mapVPNs[sIP] = EntVPNInfo{
			IsVPN:		false,
			IsInCheck:	true,
			UpdatedAt:	time.Now().Unix(),
		};
		MuVPN.Unlock();
	} else {
		return;
	}

	clientHttp := http.Client{
		Timeout: 15 * time.Second,
	}
	respHttp, errHttp := clientHttp.Get(fmt.Sprintf("http://check.getipintel.net/check.php?ip=%s&contact=%s&format=json&flags=f", sIP, settings.GetIPIntelContact));
	if (errHttp != nil) {
		return;
	}
	defer respHttp.Body.Close();
	if (respHttp.StatusCode != 200) {
		return;
	}
	byRespBody, errRespBody := io.ReadAll(respHttp.Body);
	if (errRespBody != nil) {
		return;
	}

	sStatus, errStatus := jsonparser.GetString(byRespBody, "status");
	if (errStatus != nil || sStatus != "success") {
		return;
	}
	f64Result, errResult := jsonparser.GetFloat(byRespBody, "result");
	if (errResult != nil || f64Result < 0.0 || f64Result > 1.0) {
		return;
	}
	MuVPN.Lock();
	if (f64Result >= 0.8) {
		mapVPNs[sIP] = EntVPNInfo{
			IsVPN:		true,
			IsInCheck:	false,
			UpdatedAt:	time.Now().Unix(),
		};
		fmt.Printf("Received IP VPN info (%s): %.09f (is VPN)\n", sIP, f64Result);
	} else {
		mapVPNs[sIP] = EntVPNInfo{
			IsVPN:		false,
			IsInCheck:	false,
			UpdatedAt:	time.Now().Unix(),
		};
		fmt.Printf("Received IP VPN info (%s): %.09f (is not VPN)\n", sIP, f64Result);
	}
	MuVPN.Unlock();
}

func IsVPN(sIP string) bool { //thread safe, fast
	MuVPN.RLock();
	oVPNInfo, bFound := mapVPNs[sIP];
	if (bFound && oVPNInfo.IsVPN) {
		MuVPN.RUnlock();
		return true;
	}
	MuVPN.RUnlock();
	return false;
}
