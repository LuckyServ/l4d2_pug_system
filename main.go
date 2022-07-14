package main

import (
	"fmt"
	"./settings"
	"./players"
	"./database"
	"./api"
	"./players/auth"
	//"time"
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

	//HTTP server init
	api.GinInit();

	go players.WatchOnline();

	//Block until shutdown command is received
	fmt.Printf("End: %v\n", <-api.ChShutdown);
}
