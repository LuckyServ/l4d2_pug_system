package main

import (
	"fmt"
	"./settings"
	"./players"
	"./database"
	"./players/auth"
	//"time"
)

var chShutdown chan bool = make(chan bool);


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
	ginInit();

	go players.WatchOnline();

	//Block until shutdown command is received
	fmt.Printf("End: %v\n", <-chShutdown);
}

func PerformShutDown() {
	chShutdown <- true;
}
