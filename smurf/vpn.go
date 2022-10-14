package smurf

import (
	"time"
	"sync"
	"io"
	"github.com/buger/jsonparser"
	"net/http"
	"fmt"
	"strconv"
	"../settings"
	"../database"
)

type EntVPNInfo struct {
	IsVPN			bool
	IsInCheck		bool
	UpdatedAt		int64 //unix time in seconds
}

var MapVPNs = make(map[string]EntVPNInfo);
var MuVPN sync.RWMutex;

var chVPNReqAdd = make(chan bool);
var chNewVPNReqAllowed = make(chan bool);
var arGetIPIntelRequests []time.Time;
var i64Minute int64 = 60*1000000000;
var i64Day int64 = 24*60*60*1000000000;


func Watchers() {
	go HandleLimits();
	go ClearOutdated();
}

func HandleLimits() {
	for {
		select {
		case chNewVPNReqAllowed <- func()(bool) {
			iCount := len(arGetIPIntelRequests);
			oNow := time.Now();
			iCutAt := -1;
			iPerDay, iPerMin := 0, 0;
			for i := iCount - 1; i >= 0; i-- {
				i64TimeDiff := int64(oNow.Sub(arGetIPIntelRequests[i]));
				if (i64TimeDiff <= i64Minute) {
					iPerMin++;
				}
				if (i64TimeDiff <= i64Day) {
					iPerDay++;
				} else {
					iCutAt = i;
					break;
				}
			}
			if (iCutAt >= 0) {
				arGetIPIntelRequests = arGetIPIntelRequests[(iCutAt + 1):];
			}
			if (iPerMin < 10 && iPerDay < 300) {
				arGetIPIntelRequests = append(arGetIPIntelRequests, time.Now());
				return true;
			}
			return false;
		}():
		}

	}
}

func ClearOutdated() {
	for {
		time.Sleep(86400 * time.Second); //24 hours
		var arRemoveIP []string;
		MuVPN.RLock();
		i64CurTime := time.Now().Unix();
		for sIP, oVPNInfo := range MapVPNs {
			if (oVPNInfo.UpdatedAt + 1209600/*2weeks*/ <= i64CurTime) {
				arRemoveIP = append(arRemoveIP, sIP);
			}
		}
		MuVPN.RUnlock();
		MuVPN.Lock();
		for _, sIP := range arRemoveIP {
			delete(MapVPNs, sIP);
		}
		MuVPN.Unlock();
	}
}

func CheckVPN(sIP string) {
	MuVPN.RLock();
	bShouldCheck := true;
	oVPNInfo, bFound := MapVPNs[sIP];
	if (bFound) {
		i64CurTimeSec := time.Now().Unix();
		if (oVPNInfo.IsInCheck && oVPNInfo.UpdatedAt + 60/*1min*/ > i64CurTimeSec) {
			bShouldCheck = false;
		} else if (!oVPNInfo.IsInCheck && oVPNInfo.UpdatedAt + 604800/*1week*/ > i64CurTimeSec) {
			bShouldCheck = false;
		}
	}
	MuVPN.RUnlock();

	if (bShouldCheck && <-chNewVPNReqAllowed) {
		MuVPN.Lock();
		MapVPNs[sIP] = EntVPNInfo{
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
	sResult, errResult := jsonparser.GetString(byRespBody, "result");
	f64Result, errParseFloat := strconv.ParseFloat(sResult, 64);
	if (errResult != nil || errParseFloat != nil || f64Result < 0.0 || f64Result > 1.0) {
		return;
	}
	MuVPN.Lock();
	var oNewVPNInfo EntVPNInfo;
	if (f64Result >= 0.8) {
		oNewVPNInfo = EntVPNInfo{
			IsVPN:		true,
			IsInCheck:	false,
			UpdatedAt:	time.Now().Unix(),
		};
	} else {
		oNewVPNInfo = EntVPNInfo{
			IsVPN:		false,
			IsInCheck:	false,
			UpdatedAt:	time.Now().Unix(),
		};
	}
	MapVPNs[sIP] = oNewVPNInfo;
	go database.SaveVPNInfo(database.DatabaseVPNInfo{
		IsVPN:			oNewVPNInfo.IsVPN,
		IP:				sIP,
		UpdatedAt:		oNewVPNInfo.UpdatedAt,
	});
	MuVPN.Unlock();
}

func IsVPN(sIP string) bool { //true if checked and vpn, false otherwise
	MuVPN.RLock();
	oVPNInfo, bFound := MapVPNs[sIP];
	if (bFound && oVPNInfo.IsVPN) {
		MuVPN.RUnlock();
		return true;
	}
	MuVPN.RUnlock();
	return false;
}

func IsNotVPN(sIP string) bool { //true if checked and not vpn, false otherwise
	MuVPN.RLock();
	oVPNInfo, bFound := MapVPNs[sIP];
	if (bFound && !oVPNInfo.IsInCheck && !oVPNInfo.IsVPN) {
		MuVPN.RUnlock();
		return true;
	}
	MuVPN.RUnlock();
	return false;
}

func RestoreVPNInfo() bool { //no need to lock anything
	arDBVpnInfo := database.RestoreVPNInfo();
	for _, oDBVpnInfo := range arDBVpnInfo {
		MapVPNs[oDBVpnInfo.IP] = EntVPNInfo{
			IsVPN:			oDBVpnInfo.IsVPN,
			IsInCheck:		false,
			UpdatedAt:		oDBVpnInfo.UpdatedAt,
		};
	}
	return true;
}