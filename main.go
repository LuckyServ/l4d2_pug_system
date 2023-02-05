package main

import (
	"fmt"
	"./settings"
	"./players"
	"./database"
	"./api"
	"./rating"
	"./queue"
	"./games"
	"./chat"
	"./smurf"
	"./bans"
	"./players/auth"
	"time"
	"./utils"
	"./streams"
	//"crypto/rand"
	//"math/big"
	//"encoding/json"
	//"encoding/base64"
)

/*
Mutex lock order:
MuSessions
MuPlayers
MuGames
MuDatabase
MuStreams
MuVPN
MuAuth
MuAuthTwitch
MuA2S
*/


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
	if (!bans.RestoreBans()) {
		return;
	}
	if (!auth.RestoreSessions()) {
		return;
	}
	i64CurTime := time.Now().UnixMilli();
	players.I64LastPlayerlistUpdate = i64CurTime;
	rating.GeneratePairingVariants();
	api.SetupOpenID();

	//HTTP server init
	api.GinInit();

	go players.Watchers();
	go auth.Watchers();
	go api.AuthRatelimits();
	go queue.Watchers();
	go games.Watchers();
	go chat.ChannelWatchers();
	go smurf.Watchers();
	go bans.Watchers();
	go utils.Watchers();
	go rating.Watchers();
	go streams.Watchers();


	//Watch memory leaks
	go WatchMemory();


	//Test
	//go TestingPlayersFromMain();
	//go TestingCreateBansFromMain();


	//Block until shutdown command is received
	fmt.Printf("End: %v\n", <-api.ChShutdown);
}

/*func TestingCreateBansFromMain() {
	sBanReason := base64.StdEncoding.EncodeToString([]byte("Test ban test ban test ban"));
	for i := 1; i <= 554; i++ {
		sGenSteamID64, _ := utils.GenerateRandomString(17, "12345689");
		sGenName, _ := utils.GenerateRandomString(10, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz");
		oBanRecord := bans.EntBanRecord{
			NicknameBase64:		base64.StdEncoding.EncodeToString([]byte(sGenName)),
			SteamID64:			sGenSteamID64,
			BannedBySteamID64:	"auto",
			CreatedAt:			time.Now().UnixMilli(),
			BanLength:			63072000000,
			BanReasonBase64:	sBanReason,
		};
		bans.ArrayBanRecords = append(bans.ArrayBanRecords, oBanRecord);
		time.Sleep(2 * time.Millisecond);
	}
}*/

/*func TestingPlayersFromMain() {
	time.Sleep(10 * time.Second);


	for i := 1; i <= 7; i++ {
		sGenSteamID64, _ := utils.GenerateRandomString(17, "12345689");
		sGenName, _ := utils.GenerateRandomString(10, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz");
		pPlayer := &players.EntPlayer{
			SteamID64:			sGenSteamID64,
			NicknameBase64:		base64.StdEncoding.EncodeToString([]byte(sGenName)),
			Mmr:				int((time.Now().UnixNano() % 3000) + 1),
			MmrUncertainty:		settings.DefaultMmrUncertainty,
			ProfValidated:		true,
			RulesAccepted:		true,
			LastActivity:		time.Now().UnixMilli(),
			IsOnline:			true,
			OnlineSince:		time.Now().UnixMilli(),
			LastLobbyActivity:	time.Now().UnixMilli(), //Last lobby activity //unix timestamp in milliseconds
		};
		players.MapPlayers[sGenSteamID64] = pPlayer;
		players.ArrayPlayers = append(players.ArrayPlayers, pPlayer);
	}
	players.I64LastPlayerlistUpdate = time.Now().UnixMilli();


	for _, pPlayer := range players.ArrayPlayers {
		if (pPlayer.IsOnline && !pPlayer.IsInLobby && !pPlayer.IsInGame) {
			lobby.JoinAny(pPlayer);
		}
	}

	time.Sleep(5 * time.Second);

	for _, pPlayer := range players.ArrayPlayers {
		if (pPlayer.IsInLobby && !pPlayer.IsReadyInLobby) {
			//time.Sleep(1 * time.Second);
			lobby.Ready(pPlayer);
		}
	}

}*/