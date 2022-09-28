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

func SelectBestAvailableServer(arPlayers []*players.EntPlayer) string { //must unlock players inside, but they are locked outside

	//find player ping weight based on their region and time of the day
	for _, pPlayer := range arPlayers {
		if (len(pPlayer.GameServerPings) == len(settings.GameServers)) {
			var sLowestPingIP string;
			var iLowestPing int = 350;
			for sIP, iPing := range pPlayer.GameServerPings {
				if (iPing < iLowestPing) {
					iLowestPing = iPing;
					sLowestPingIP = sIP;
				}
			}
			sRegion := GetRegionFromIP(sLowestPingIP);
			iCurHour := time.Now().Hour();
			if (iCurHour >= 0 && iCurHour < 8 && sRegion == "america") { //America time
				pPlayer.GameServerPingWeight = 2;
			} else if (iCurHour >= 8 && iCurHour < 16 && sRegion == "asia") { //Asia time
				pPlayer.GameServerPingWeight = 2;
			} else if (iCurHour >= 16 && iCurHour < 24 && sRegion == "europe") { //Europe time
				pPlayer.GameServerPingWeight = 2;
			} else {
				pPlayer.GameServerPingWeight = 1;
			}
		} else {
			pPlayer.GameServerPingWeight = 0;
		}
	}


	//select server based on weighted average ping
	var iBestWeightedPing int = 350;
	oBestWeighted := settings.GameServers[0];
	iSumOfWeights := SumOfWeights(arPlayers);
	if (iSumOfWeights > 0) {
		for _, oServer := range settings.GameServers {
			iWeightedAverage := SumWeightedPings(arPlayers, oServer.IP) / iSumOfWeights;
			if (iWeightedAverage < iBestWeightedPing) {
				iBestWeightedPing = iWeightedAverage;
				oBestWeighted = oServer;
			}
		}
	}


	//select server based on max ping within region selected above
	var arGameServers []string;
	var arMaxPing []int;
	var arGameServersWW []string;
	var arMaxPingWW []int;

	for _, oServer := range settings.GameServers {
		var iMaxPing int;
		for _, pPlayer := range arPlayers {
			if (pPlayer.GameServerPings[oServer.IP] > iMaxPing) {
				iMaxPing = pPlayer.GameServerPings[oServer.IP];
			}
		}
		if (iMaxPing == 0) {
			iMaxPing = 350;
		}
		for _, sPort := range oServer.Ports {
			arGameServersWW = append(arGameServersWW, oServer.IP+":"+sPort);
			arMaxPingWW = append(arMaxPingWW, iMaxPing);
			if (oServer.Region == oBestWeighted.Region) {
				arGameServers = append(arGameServers, oServer.IP+":"+sPort);
				arMaxPing = append(arMaxPing, iMaxPing);
			}
		}
	}

	players.MuPlayers.Unlock();

	//sort by maxping, exclude occupied and outdated servers, return 1st server
	sChosenServer := GetAvailableServer(arGameServers, arMaxPing);


	//if no servers available, search again within all regions
	if (sChosenServer == "") {
		sChosenServer = GetAvailableServer(arGameServersWW, arMaxPingWW);
	}

	//return
	return sChosenServer;
}


func GetAvailableServer(arGameServers []string, arMaxPing []int) string {
	SortByMaxPing(arGameServers, arMaxPing);
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

func SumWeightedPings(arPlayers []*players.EntPlayer, sIP string) int {
	var iSum int;
	for _, pPlayer := range arPlayers {
		iPing, _ := pPlayer.GameServerPings[sIP];
		iSum = iSum + (iPing * pPlayer.GameServerPingWeight);
	}
	return iSum;
}

func SumOfWeights(arPlayers []*players.EntPlayer) int {
	var iSum int;
	for _, pPlayer := range arPlayers {
		iSum = iSum + pPlayer.GameServerPingWeight;
	}
	return iSum;
}

func SortByMaxPing(arGameServers []string, arMaxPing []int) {
	iSize := len(arGameServers);
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
}

func GetRegionFromIP(sIP string) string {
	for _, oServer := range settings.GameServers {
		if (oServer.IP == sIP) {
			return oServer.Region;
		}
	}
	return ""; //can't happen
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