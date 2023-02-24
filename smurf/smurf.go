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

type IPInfo struct {
	IP				string			`json:"ip"`
	FirstUsed		int64			`json:"first_used"`
	GeoIPDate		int64			`json:"geoip_date"`
	Location		string			`json:"location"`
	ISP				string			`json:"isp"`
	IsVPN			int				`json:"is_vpn"`
}


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

func GetIPInfo(sSteamID64 string) []IPInfo {
	var arIPInfo []IPInfo;

	clientHttp := http.Client{
		Timeout: 15 * time.Second,
	}
	respHttp, errHttp := clientHttp.Get(fmt.Sprintf("%s/getusedips?auth_key=%s&steamid64=%s", settings.SmurfHost, settings.SmurfAuthKey, sSteamID64));
	if (errHttp == nil) {
		if (respHttp.StatusCode == 200) {
			byRespBody, errRespBody := io.ReadAll(respHttp.Body);
			if (errRespBody == nil) {
				bSuccess, _ := jsonparser.GetBoolean(byRespBody, "success");
				if (bSuccess) {
					jsonparser.ArrayEach(byRespBody, func(valueIP []byte, dataType jsonparser.ValueType, offset int, err error) {


						sIP, _ := jsonparser.GetString(valueIP, "ip");
						i64FirstUsed, _ := jsonparser.GetInt(valueIP, "first_used");
						i64GeoIPDate, _ := jsonparser.GetInt(valueIP, "geoip_date");
						sLocation, _ := jsonparser.GetString(valueIP, "location");
						sISP, _ := jsonparser.GetString(valueIP, "isp");
						i64IsVPN, _ := jsonparser.GetInt(valueIP, "is_vpn");

						arIPInfo = append(arIPInfo, IPInfo{
							IP:				sIP,
							FirstUsed:		i64FirstUsed,
							GeoIPDate:		i64GeoIPDate,
							Location:		sLocation,
							ISP:			sISP,
							IsVPN:			int(i64IsVPN),
						});


					}, "ips");
				}
			}
		}
		respHttp.Body.Close();
	}

	return arIPInfo;
}