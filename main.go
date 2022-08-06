package main

import (
	"fmt"
	"./settings"
	"./players"
	"./database"
	"./lobby"
	"./api"
	"./games"
	"./players/auth"
	"time"
	//"./utils"
	//"crypto/rand"
	//"math/big"
	//"encoding/json"
)


func main() {
	fmt.Printf("Started.\n");

	//Parse settings
	if (!settings.Parse()) {
		return;
	}

	//Permanent global database connection
	if (!database.DatabaseConnect()) {
		return;
	}

	//Restore from database
	if (!players.RestorePlayers()) {
		return;
	}
	if (!auth.RestoreSessions()) {
		return;
	}
	lobby.I64LastLobbyListUpdate = time.Now().UnixMilli();

	//HTTP server init
	api.GinInit();

	go players.WatchOnline(); //send players to offline mode
	go api.AuthRatelimits(); //limit authorization requests per ip
	go lobby.WatchLobbies(); //watch lobbies for ready and timeout
	go games.ChannelWatchers(); //watch various game related channels




	//Test pinging
	//go PingTestingFromMain();


	//Block until shutdown command is received
	fmt.Printf("End: %v\n", <-api.ChShutdown);
}

/*func PingTestingFromMain() {
	for _, pPlayer := range players.ArrayPlayers {
		pPlayer.GameServerPings = make(map[string]int, len(settings.HardwareServers));
	}
	for i := 1; i <= 10; i++ {
		fmt.Printf("\n");
		for _, pPlayer := range players.ArrayPlayers {
			fmt.Printf("%v\n", pPlayer.GameServerPings);
		}
		time.Sleep(1 * time.Second);
	}

	
	iTryCount := 0;
	for {
		fmt.Printf("Retrieving serverlist\n");
		arAvailGameSrvs := games.GetAvailableServers();
		fmt.Printf("Serverlist: %v\n", arAvailGameSrvs);

		players.MuPlayers.Lock();
		sIPPORT := games.SelectBestAvailableServer(players.ArrayPlayers, arAvailGameSrvs);
		fmt.Printf("sIPPORT == \"%s\"\n", sIPPORT);

		games.MuGames.Lock();
		bSuccess := (sIPPORT != "");

		if (bSuccess) {
			games.MuGames.Unlock();
			players.MuPlayers.Unlock();
			fmt.Printf("Success\n");
			break;
		} else {

			games.MuGames.Unlock();
			players.MuPlayers.Unlock();

			iTryCount++;
			if (iTryCount >= settings.AvailGameSrvsMaxTries) {
				games.MuGames.Lock();
				players.MuPlayers.Lock();
				games.MuGames.Unlock();
				players.MuPlayers.Unlock();
				return;
			}
		}
		time.Sleep(60 * time.Second);
	}
}*/