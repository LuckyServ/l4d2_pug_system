package main

import (
	"fmt"
	//"database/sql"
	//_ "github.com/lib/pq"
	"./settings"
)

var bStateShutdown bool;
var chShutdown chan bool = make(chan bool);

/*var dbConn *sql.DB;
var dbErr error;
var sDbConn_string = "postgres://user:password@address:port/database";*/ //PostgreSQL connection credentials


func main() {
	fmt.Printf("Started.\n");

	//Parse settings
	if (!settings.Parse()) {
		return;
	}

	//Permanent global database connection
	/*if (!DatabaseConnect()) {
		return;
	}*/

	//HTTP server init
	ginInit();


	//Block until shutdown command is received
	fmt.Printf("End: %v\n", <-chShutdown);
}


/*func DatabaseConnect() bool {
	dbConn, dbErr = sql.Open("postgres", sDbConn_string);
	if (dbErr != nil) {
		fmt.Printf("Error connecting to database: Open\n");
		return false;
	}
	dbErr = dbConn.Ping();
	if (dbErr != nil) {
		fmt.Printf("Error connecting to database: Ping\n");
		return false;
	}
	fmt.Printf("Database connection successfull\n");
	return true;
}*/

