package games

import (
	"fmt"
	"../players"
	"../settings"
	"../rating"
	"time"
)


func Control(pGame *EntGame) {

	//Create
	MuGames.Lock();
	players.MuPlayers.Lock();

	Create(pGame);
	pGame.State = StateCreated;
	SetLastUpdated(pGame.PlayersUnpaired);

	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Choose maps
	i64CmpgnIdx := time.Now().UnixNano() % int64(len(settings.CampaignNames));
	MuGames.Lock();
	players.MuPlayers.Lock();
	pGame.CampaignName = settings.CampaignNames[i64CmpgnIdx];
	pGame.Maps = settings.MapPool[i64CmpgnIdx];
	pGame.State = StateCampaignChosen;
	SetLastUpdated(pGame.PlayersUnpaired);
	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Pair players
	MuGames.Lock();
	players.MuPlayers.Lock();

	pGame.PlayersA, pGame.PlayersB = rating.Pair(pGame.PlayersUnpaired);

	pGame.State = StateTeamsPicked;
	SetLastUpdated(pGame.PlayersUnpaired);
	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Request pings
	MuGames.Lock();
	players.MuPlayers.Lock();
	for _, pPlayer := range pGame.PlayersUnpaired {
		pPlayer.GameServerPings = make(map[string]int, len(settings.HardwareServers));
	}
	pGame.State = StateWaitPings;
	SetLastUpdated(pGame.PlayersUnpaired);
	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Wait for pings
	time.Sleep(time.Duration(settings.MaxPingWait) * time.Second);


	//Cancel the ping request
	MuGames.Lock();
	players.MuPlayers.Lock();
	pGame.State = StateSelectServer;
	SetLastUpdated(pGame.PlayersUnpaired);
	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Select best available server based on pings and availability (a2s requests here)
	iTryCount := 0;
	for {
		arAvailGameSrvs := GetAvailableServers(); //long execution (a2s); no need to lock anything

		players.MuPlayers.Lock();
		sIPPORT := SelectBestAvailableServer(pGame.PlayersUnpaired, arAvailGameSrvs);

		MuGames.Lock();
		bSuccess := (sIPPORT != "");
		if (bSuccess) {
			for _, pGameI := range ArrayGames {
				if (pGameI.ServerIP == sIPPORT) {
					bSuccess = false;
					break;
				}
			}
		}

		if (bSuccess) {
			pGame.ServerIP = sIPPORT;
			pGame.State = StateWaitPlayersJoin;
			SetLastUpdated(pGame.PlayersUnpaired);
			MuGames.Unlock();
			players.MuPlayers.Unlock();
			break;
		} else {

			pGame.State = StateNoServers;
			SetLastUpdated(pGame.PlayersUnpaired);
			MuGames.Unlock();
			players.MuPlayers.Unlock();

			iTryCount++;
			if (iTryCount >= settings.AvailGameSrvsMaxTries) { //destroy lobby if too many tries
				MuGames.Lock();
				players.MuPlayers.Lock();
				Destroy(pGame);
				MuGames.Unlock();
				players.MuPlayers.Unlock();
				return;
			}
		}
		time.Sleep(60 * time.Second); //check once per minute
	}


	chRUpExpired := make(chan bool);
	//Wait for first readyup
	go func(chRUpExpired chan bool)() {
		time.Sleep(time.Duration(settings.FirstReadyUpExpire) * time.Second);
		select {
		case chRUpExpired <- true:
		default:
		}
	}(chRUpExpired);

	chFullRUpReceived := make(chan bool);
	MuGames.Lock();
	pGame.ReceiverFullRUP = chFullRUpReceived;
	MuGames.Unlock();
	select {
	case <-chRUpExpired:
		MuGames.Lock();
		players.MuPlayers.Lock();
		pGame.State = StateReadyUpExpired;
		SetLastUpdated(pGame.PlayersUnpaired);
		MuGames.Unlock();
		players.MuPlayers.Unlock();
	case <-chFullRUpReceived:
		MuGames.Lock();
		players.MuPlayers.Lock();
		pGame.State = StateGameProceeds;
		SetLastUpdated(pGame.PlayersUnpaired);
		MuGames.Unlock();
		players.MuPlayers.Unlock();
	}


	fmt.Printf("Game proceeds\n");
	select{};
	//Game proceeds
	//Game ended, settle results
	//Destroy Game
}