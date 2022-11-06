package games

import (
	//"fmt"
	"../settings"
	"../players"
	"../utils"
	"github.com/rumblefrog/go-a2s"
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

func SelectBestAvailableServer(arPlayersA []*players.EntPlayer, arPlayersB []*players.EntPlayer) string { //must unlock players inside, but they are locked outside


	var arGameServers []string;
	var arPriority []int;

	for _, oServer := range settings.GameServers {
		var iTeamPingDiff, iAvgPing = CalcPings(arPlayersA, arPlayersB, oServer.IP);
		if (iAvgPing <= 200 && iTeamPingDiff <= 80) { //limits for playable conditions
			for _, sPort := range oServer.Ports {
				arGameServers = append(arGameServers, oServer.IP+":"+sPort);
				arPriority = append(arPriority, (iTeamPingDiff + iAvgPing) / 2);
			}
		}
	}


	players.MuPlayers.Unlock();


	return GetAvailableServer(arGameServers, arPriority);
}

func CalcPings(arPlayersA []*players.EntPlayer, arPlayersB []*players.EntPlayer, sIP string) (int, int) {
	var iTeamPingDiff int;
	var iSumOfPings, iNumOfPings, iAvgPing int;
	var iSumOfPingsA, iNumOfPingsA, iAvgPingA int;
	for _, pPlayer := range arPlayersA {
		iPing, bPingExists := pPlayer.GameServerPings[sIP];
		if (bPingExists && iPing > 0) {
			iSumOfPingsA = iSumOfPingsA + iPing;
			iNumOfPingsA++;
			iSumOfPings = iSumOfPings + iPing;
			iNumOfPings++;
		}
	}
	if (iNumOfPingsA > 0) {
		iAvgPingA = iSumOfPingsA / iNumOfPingsA;
	}

	var iSumOfPingsB, iNumOfPingsB, iAvgPingB int;
	for _, pPlayer := range arPlayersB {
		iPing, bPingExists := pPlayer.GameServerPings[sIP];
		if (bPingExists && iPing > 0) {
			iSumOfPingsB = iSumOfPingsB + iPing;
			iNumOfPingsB++;
			iSumOfPings = iSumOfPings + iPing;
			iNumOfPings++;
		}
	}
	if (iNumOfPingsB > 0) {
		iAvgPingB = iSumOfPingsB / iNumOfPingsB;
	}

	if (iAvgPingA > 0 && iAvgPingB > 0) {
		iTeamPingDiff = iAvgPingA - iAvgPingB;
		if (iTeamPingDiff < 0) {iTeamPingDiff = iTeamPingDiff * -1;};
	}


	if (iNumOfPings > 0) {
		iAvgPing = iSumOfPings / iNumOfPings;
	}

	return iTeamPingDiff, iAvgPing;
}


func GetAvailableServer(arGameServers []string, arPriority []int) string {
	SortByPriority(arGameServers, arPriority);
	var arEmptyGameSrvs []string;
	var arQueryCh []chan int;
	for range arGameServers {
		arQueryCh = append(arQueryCh, make(chan int));
	}
	for i, sIPPORT := range arGameServers {
		go GetPlayersCount(arQueryCh[i], sIPPORT);
	}
	for i, sIPPORT := range arGameServers {
		iCount := <-arQueryCh[i];
		if (iCount == 0) {
			arEmptyGameSrvs = append(arEmptyGameSrvs, sIPPORT);
		}
	}
	if (len(arEmptyGameSrvs) > 0) {
		return arEmptyGameSrvs[0];
	}
	return "";
}

func SortByPriority(arGameServers []string, arPriority []int) {
	iSize := len(arGameServers);
	if (iSize > 1) {
		bSorted := false;
		for !bSorted {
			bSorted = true;
			for i := 1; i < iSize; i++ {
				if (arPriority[i] < arPriority[i - 1]) {
					arPriority[i], arPriority[i - 1] = arPriority[i - 1], arPriority[i]; //switch
					arGameServers[i], arGameServers[i - 1] = arGameServers[i - 1], arGameServers[i]; //switch
					if (bSorted) {
						bSorted = false;
					}
				}
			}
			if (!bSorted) {
				for i := iSize - 2; i >= 0; i-- {
					if (arPriority[i] > arPriority[i + 1]) {
						arPriority[i], arPriority[i + 1] = arPriority[i + 1], arPriority[i]; //switch
						arGameServers[i], arGameServers[i + 1] = arGameServers[i + 1], arGameServers[i]; //switch
					}
				}
			}
		}
	}
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