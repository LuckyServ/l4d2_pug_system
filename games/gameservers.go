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
			var iLowestPing int = 999999;
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


	//select region based on weighted average ping
	var iBestWeightedPing int = 350;
	var sSelectedRegion string = "europe";
	iSumOfWeights := SumOfWeights(arPlayers);
	if (iSumOfWeights > 0) {
		for _, oServer := range settings.GameServers {
			iWeightedAverage := SumWeightedPings(arPlayers, oServer.IP) / iSumOfWeights;
			if (iWeightedAverage < iBestWeightedPing) {
				iBestWeightedPing = iWeightedAverage;
				sSelectedRegion = oServer.Region;
			}
		}
	}


	//get servers within region selected above
	var arGameServers []string;
	var arWAvgPing []int;

	for _, oServer := range settings.GameServers {
		var iWAvgPing int;
		if (iSumOfWeights > 0) {
			iWAvgPing = SumWeightedPings(arPlayers, oServer.IP) / iSumOfWeights;
		}
		if (iWAvgPing == 0) {
			iWAvgPing = 999999;
		}
		for _, sPort := range oServer.Ports {
			if (oServer.Region == sSelectedRegion) {
				arGameServers = append(arGameServers, oServer.IP+":"+sPort);
				arWAvgPing = append(arWAvgPing, iWAvgPing);
			}
		}
	}

	players.MuPlayers.Unlock();

	//sort by weighted average ping, exclude occupied and outdated servers, return 1st server
	return GetAvailableServer(arGameServers, arWAvgPing);
}


func GetAvailableServer(arGameServers []string, arWAvgPing []int) string {
	SortByWAvgPing(arGameServers, arWAvgPing);
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

func SortByWAvgPing(arGameServers []string, arWAvgPing []int) {
	iSize := len(arGameServers);
	if (iSize > 1) {
		bSorted := false;
		for !bSorted {
			bSorted = true;
			for i := 1; i < iSize; i++ {
				if (arWAvgPing[i] < arWAvgPing[i - 1]) {
					arWAvgPing[i], arWAvgPing[i - 1] = arWAvgPing[i - 1], arWAvgPing[i]; //switch
					arGameServers[i], arGameServers[i - 1] = arGameServers[i - 1], arGameServers[i]; //switch
					if (bSorted) {
						bSorted = false;
					}
				}
			}
			if (!bSorted) {
				for i := iSize - 2; i >= 0; i-- {
					if (arWAvgPing[i] > arWAvgPing[i + 1]) {
						arWAvgPing[i], arWAvgPing[i + 1] = arWAvgPing[i + 1], arWAvgPing[i]; //switch
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