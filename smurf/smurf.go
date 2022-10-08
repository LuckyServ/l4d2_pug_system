package smurf

import (
	"fmt"
	"../settings"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func AnnounceIP(sIP string) {
	go CheckVPN(sIP);
}

func AnnounceIPAndKey(sSteamID64 string, sIP string, sNickname string, sUniqueKey string) { //thread safe, fast
	go CheckVPN(sIP);

	bKeyValid, _ := regexp.MatchString(`^[0-9a-z]{16,200}$`, sUniqueKey);
	if (bKeyValid) {
		go CheckIPAndKey(sSteamID64, sIP, sNickname, sUniqueKey);
	} else {
		//ban
	}
}

func CheckIPAndKey(sSteamID64 string, sIP string, sNickname string, sUniqueKey string) {
	clientHttp := http.Client{
		Timeout: 15 * time.Second,
	}
	respHttp, errHttp := clientHttp.Get(fmt.Sprintf("%s/logipandkey?auth_key=%s&steamid64=%s&unique_key=%s&nickname=%s&ip=%s", settings.SmurfHost, settings.SmurfAuthKey, sSteamID64, sUniqueKey, url.QueryEscape(sNickname), sIP));
	if (errHttp == nil) {
		respHttp.Body.Close();
	}
}