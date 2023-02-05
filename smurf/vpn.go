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

var MapVPNs = make(map[string]EntVPNInfo);
var MuVPN sync.RWMutex;

var chVPNReqAdd = make(chan bool);
var chNewVPNReqAllowed = make(chan bool);
var arGetIPIntelRequests []time.Time;
var i64Minute int64 = 60*1000000000;
var i64Day int64 = 24*60*60*1000000000;


func Watchers() {
	go ClearOutdated();
}

func ClearOutdated() {
	for {
		time.Sleep(86400 * time.Second); //24 hours
		var arRemoveIP []string;
		MuVPN.RLock();
		i64CurTime := time.Now().Unix();
		for sIP, oVPNInfo := range MapVPNs {
			if (oVPNInfo.UpdatedAt + 604800/*1week*/ <= i64CurTime) {
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
		if (oVPNInfo.IsInCheck && oVPNInfo.UpdatedAt + 20/*20sec*/ > i64CurTimeSec) {
			bShouldCheck = false;
		} else if (!oVPNInfo.IsInCheck && oVPNInfo.UpdatedAt + 604800/*1week*/ > i64CurTimeSec) {
			bShouldCheck = false;
		}
	}
	MuVPN.RUnlock();

	if (bShouldCheck) {
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
	respHttp, errHttp := clientHttp.Get(fmt.Sprintf("%s/isvpn?auth_key=%s&ip=%s", settings.SmurfHost, settings.SmurfAuthKey, sIP));
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

	bSuccess, errSuccess := jsonparser.GetBoolean(byRespBody, "success");
	if (errSuccess != nil || bSuccess != true) {
		return;
	}
	i64Result, errResult := jsonparser.GetInt(byRespBody, "isvpn");
	if (errResult != nil || i64Result == 0) {
		return;
	}
	iResult := int(i64Result);
	MuVPN.Lock();
	var oNewVPNInfo EntVPNInfo;
	if (iResult == 2) {
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
