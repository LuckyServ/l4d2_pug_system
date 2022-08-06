package games

import (
	//"fmt"
	"../settings"
	"../players"
	"github.com/rumblefrog/go-a2s"
	"strings"
	"time"
)


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
					bSorted = false;
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

func GetPlayersCount(chCount chan int, sIPPORT string) {
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
	chCount <- int(vInfo.Players);
	return;
}