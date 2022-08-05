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
	for {
		fmt.Printf("\n");
		for _, pPlayer := range players.ArrayPlayers {
			fmt.Printf("%v\n", pPlayer.GameServerPings);
		}
		time.Sleep(1 * time.Second);
	}
}*/