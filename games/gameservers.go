package games

import (
	//"fmt"
	"../settings"
	"../players"
	"../utils"
	"github.com/rumblefrog/go-a2s"
	"strings"
	"time"
	"io/ioutil"
	"net/http"
	"github.com/buger/jsonparser"
	"strconv"
)

var sLatestGameVersion string;


func CheckVersion() {
	for {
		client := http.Client{
			Timeout: 10 * time.Second,
		}
		resHttp, errHttp := client.Get("https://api.steampowered.com/ISteamApps/UpToDateCheck/v1?appid=550&version=0&format=json");
		if (errHttp == nil) {
			if (resHttp.StatusCode == 200) {
				byResBody, errResBody := ioutil.ReadAll(resHttp.Body);
				if (errResBody == nil) {
					i64BufferVersion, errVersion := jsonparser.GetInt(byResBody, "response", "required_version");
					if (errVersion == nil) {
						sBufferVersion := strconv.FormatInt(i64BufferVersion, 10);
						sLatestGameVersion = utils.InsertDots(sBufferVersion, 1);
					}
				}
			}
			resHttp.Body.Close();
		}
		time.Sleep(60 * time.Second);
	}
}

func SelectBestAvailableServer(arPlayers []*players.EntPlayer, arGameServersUnsorted []string) string { //Players must be locked outside
	
	var arGameServers []string;
	var arMaxPing []int;

	for _, sGameServer := range arGameServersUnsorted {
		var iMaxPing int = 0;
		for _, pPlayer := range arPlayers {
			sIP := strings.Split(sGameServer, ":")[0];
			if (pPlayer.GameServerPings[sIP] > iMaxPing) {
				iMaxPing = pPlayer.GameServerPings[sIP];
			}
		}
		arGameServers = append(arGameServers, sGameServer);
		arMaxPing = append(arMaxPing, iMaxPing);
	}

	iSize := len(arGameServers);
	if (iSize == 0) {
		return ""; //no available servers
	}
	if (iSize > 1) {
		bSorted := false;
		for !bSorted {
			bSorted = true;
			for i := 1; i < iSize; i++ {
				if (arMaxPing[i] < arMaxPing[i - 1]) {
					arMaxPing[i], arMaxPing[i - 1] = arMaxPing[i - 1], arMaxPing[i]; //switch
					arGameServers[i], arGameServers[i - 1] = arGameServers[i - 1], arGameServers[i]; //switch
					if (bSorted) {
						bSorted = false;
					}
				}
			}
			if (!bSorted) {
				for i := iSize - 2; i >= 0; i-- {
					if (arMaxPing[i] > arMaxPing[i + 1]) {
						arMaxPing[i], arMaxPing[i + 1] = arMaxPing[i + 1], arMaxPing[i]; //switch
						arGameServers[i], arGameServers[i + 1] = arGameServers[i + 1], arGameServers[i]; //switch
					}
				}
			}
		}
	}

	return arGameServers[0];
}


func GetAvailableServers() []string {
	var arEmptyGameSrvs []string;
	var arQueryCh []chan int;
	for range settings.GameServers {
		arQueryCh = append(arQueryCh, make(chan int));
	}
	for i, sIPPORT := range settings.GameServers {
		go GetPlayersCount(arQueryCh[i], sIPPORT);
	}
	for i, sIPPORT := range settings.GameServers {
		iCount := <-arQueryCh[i];
		if (iCount == 0) {
			arEmptyGameSrvs = append(arEmptyGameSrvs, sIPPORT);
		}
	}
	return arEmptyGameSrvs;
}

func GetPlayersCount(chCount chan int, sIPPORT string) { //-1 if server version is outdated or unavailable
	vHandle, vErr1 := a2s.NewClient(sIPPORT, a2s.TimeoutOption(time.Second * 4));
	if (vErr1 != nil) {
		chCount <- -1;
		return;
	}
	vInfo, vErr2 := vHandle.QueryInfo();
	if (vErr2 != nil) {
		chCount <- -1;
		return;
	}
	vHandle.Close();
	if (vInfo.Version == sLatestGameVersion && sLatestGameVersion != "") {
		chCount <- int(vInfo.Players);
	} else {
		chCount <- -1;
	}
	return;
}