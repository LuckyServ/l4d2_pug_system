package smurf

import (
	"fmt"
	"../settings"
	"net/http"
	"net/url"
	"regexp"
	"time"
	"strings"
	"io"
	"github.com/buger/jsonparser"
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
		//ban mb? clearly malicious action
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

func GetKnownAccounts(sSteamID64 string) []string { //validate steam id string outside
	clientHttp := http.Client{
		Timeout: 15 * time.Second,
	}
	respHttp, errHttp := clientHttp.Get(fmt.Sprintf("%s/getknownsmurfs?auth_key=%s&steamid64=%s", settings.SmurfHost, settings.SmurfAuthKey, sSteamID64));
	if (errHttp == nil) {
		if (respHttp.StatusCode == 200) {
			byRespBody, errRespBody := io.ReadAll(respHttp.Body);
			if (errRespBody == nil) {
				sAccountsList, errAccountsList := jsonparser.GetString(byRespBody, "accounts");
				if (errAccountsList == nil) {
					return strings.Split(sAccountsList, ",");
				}
			}
		}
		respHttp.Body.Close();
	}
	return []string{};
}