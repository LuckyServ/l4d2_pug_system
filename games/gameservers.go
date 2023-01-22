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
	"sync"
)

var sLatestGameVersion string;
var MuA2S sync.RWMutex;

var mapA2SCache map[string]int;
var i64a2sCachedAt int64;


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

func GetGameServers(arPlayersA []*players.EntPlayer, arPlayersB []*players.EntPlayer) ([]string) {
	var arGameServers []string;
	var arPriority []int;

	players.MuPlayers.RLock();
	MuGames.RLock();
	for _, oServer := range settings.GameServers {
		var iTeamPingDiff, iAvgPing, iMaxPing, _ = CalcPings(arPlayersA, arPlayersB, oServer.IP);

		iPriority := ((iTeamPingDiff + iAvgPing + (iMaxPing / 2)) / 3) + oServer.LowerPriority;

		for _, sPort := range oServer.Ports {
			arGameServers = append(arGameServers, oServer.IP+":"+sPort);
			arPriority = append(arPriority, iPriority);
		}
	}
	MuGames.RUnlock();
	players.MuPlayers.RUnlock();

	SortByPriority(arGameServers, arPriority);
	return arGameServers;
}

func CalcPings(arPlayersA []*players.EntPlayer, arPlayersB []*players.EntPlayer, sIP string) (int, int, int, int) {
	var iTeamPingDiff, iMaxPing int;
	var iMinPing int = 2000000000;
	var iSumOfPings, iNumOfPings, iAvgPing int;
	var iSumOfPingsA, iNumOfPingsA, iAvgPingA int;
	for _, pPlayer := range arPlayersA {
		if (IsPingInfoValid(pPlayer)) {
			iPing, bPingExists := pPlayer.GameServerPings[sIP];
			if (bPingExists && iPing > 0) {
				iSumOfPingsA = iSumOfPingsA + iPing;
				iNumOfPingsA++;
				iSumOfPings = iSumOfPings + iPing;
				iNumOfPings++;
				if (iPing > iMaxPing) {
					iMaxPing = iPing;
				}
				if (iPing < iMinPing) {
					iMinPing = iPing;
				}
			}
		}
	}
	if (iNumOfPingsA > 0) {
		iAvgPingA = iSumOfPingsA / iNumOfPingsA;
	}

	var iSumOfPingsB, iNumOfPingsB, iAvgPingB int;
	for _, pPlayer := range arPlayersB {
		if (IsPingInfoValid(pPlayer)) {
			iPing, bPingExists := pPlayer.GameServerPings[sIP];
			if (bPingExists && iPing > 0) {
				iSumOfPingsB = iSumOfPingsB + iPing;
				iNumOfPingsB++;
				iSumOfPings = iSumOfPings + iPing;
				iNumOfPings++;
				if (iPing > iMaxPing) {
					iMaxPing = iPing;
				}
				if (iPing < iMinPing) {
					iMinPing = iPing;
				}
			}
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

	return iTeamPingDiff, iAvgPing, iMaxPing, iMinPing;
}


func GetEmptyServers(arGameServers []string) []string {
	var arEmptyGameSrvs []string;
	var arQueryCh []chan int;

	MuA2S.Lock();

	if (i64a2sCachedAt + 10000/*10s*/ > time.Now().UnixMilli()) {
		for _, sIPPORT := range arGameServers {
			iCount, bCached := mapA2SCache[sIPPORT];
			if (bCached && iCount == 0) {
				arEmptyGameSrvs = append(arEmptyGameSrvs, sIPPORT);
			}
		}
	} else {

		mapA2SCache = make(map[string]int, len(arGameServers));
		i64a2sCachedAt = time.Now().UnixMilli();

		for range arGameServers {
			arQueryCh = append(arQueryCh, make(chan int));
		}
		for i, sIPPORT := range arGameServers {
			go GetPlayersCount(arQueryCh[i], sIPPORT);
		}
		for i, sIPPORT := range arGameServers {
			iCount := <-arQueryCh[i];
			mapA2SCache[sIPPORT] = iCount;
			if (iCount == 0) {
				arEmptyGameSrvs = append(arEmptyGameSrvs, sIPPORT);
			}
		}

	}
	
	MuA2S.Unlock();
	return arEmptyGameSrvs;
}

func GetUnreservedServers(arGameServers []string) []string {
	var arUnreservedGameSrvs []string;
	for _, sIPPORT := range arGameServers {
		if (GetGameByIP(sIPPORT) == nil) {
			arUnreservedGameSrvs = append(arUnreservedGameSrvs, sIPPORT);
		}
	}
	return arUnreservedGameSrvs;
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
	defer vHandle.Close();
	vInfo, vErr2 := vHandle.QueryInfo();
	if (vErr2 != nil) {
		chCount <- -1;
		return;
	}
	if (vInfo.Version == sLatestGameVersion && sLatestGameVersion != "" && GameRuleExists(vHandle, "l4d_ready_enabled") == 0) {
		chCount <- int(vInfo.Players);
	} else {
		chCount <- -1;
	}
	return;
}

func GameRuleExists(vHandle *a2s.Client, sRule string) int { //-1 - error, 0 - doesnt exist, 1 - exists
	vRules, vErr2 := vHandle.QueryRules();
	if (vErr2 != nil) {
		return -1;
	}
	mapRules := vRules.Rules;
	_, bExists := mapRules[sRule];
	if (bExists) {
		return 1;
	}
	return 0;
}

func IsPingInfoValid(pPlayer *players.EntPlayer) bool {
	mapPings := pPlayer.GameServerPings;
	if (mapPings != nil) {
		for _, iPing := range mapPings {
			if (iPing >= 450) {
				return false;
			}
		}
	} else {
		return false;
	}
	return true;
}