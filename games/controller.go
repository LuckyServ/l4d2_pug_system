package games

import (
	//"fmt"
	"../players"
	"../settings"
	"../rating"
	"time"
	"../database"
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


	//Wait for first readyup
	chRUpExpired := make(chan bool);
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


	//If Readyup not received, do smth
	MuGames.Lock();
	iStateBuffer := pGame.State;
	MuGames.Unlock();
	if (iStateBuffer == StateReadyUpExpired) {

		chListOfReadyPlayersExpired := make(chan bool);
		go func(chListOfReadyPlayersExpired chan bool)() {
			time.Sleep(30 * time.Second);
			select {
			case chListOfReadyPlayersExpired <- true:
			default:
			}
		}(chListOfReadyPlayersExpired);
		chListOfReadyPlayers := make(chan []string);
		MuGames.Lock();
		pGame.ReceiverReadyList = chListOfReadyPlayers;
		MuGames.Unlock();

		select {
		case <-chListOfReadyPlayersExpired:
			//Ban no one, destroy the game
			MuGames.Lock();
			players.MuPlayers.Lock();
			Destroy(pGame);
			MuGames.Unlock();
			players.MuPlayers.Unlock();
			return;
		case arReadyPlayers := <-chListOfReadyPlayers:
			if (len(arReadyPlayers) == 1 && (arReadyPlayers[0] == "none_ready" || arReadyPlayers[0] == "all_ready")) {
				//do nothing
			} else {
				//ban those who isnt ready
			}
			//Destroy the game
			MuGames.Lock();
			players.MuPlayers.Lock();
			Destroy(pGame);
			MuGames.Unlock();
			players.MuPlayers.Unlock();
			return;
		case <-chFullRUpReceived:
			//Proceed to the next state
			MuGames.Lock();
			players.MuPlayers.Lock();
			pGame.State = StateGameProceeds;
			SetLastUpdated(pGame.PlayersUnpaired);
			MuGames.Unlock();
			players.MuPlayers.Unlock();
		}
	}



	//Game proceeds
	chResult := make(chan rating.EntGameResult);
	MuGames.Lock();
	pGame.ReceiverResult = chResult;
	pGame.GameResult = rating.EntGameResult{};
	MuGames.Unlock();

	bGameProceeds := true;
	for bGameProceeds { //continuously receive results until game ends

		chResultExpired := make(chan bool);
		go func(chResultExpired chan bool)() {
			time.Sleep(60 * time.Second);
			select {
			case chResultExpired <- true:
			default:
			}
		}(chResultExpired);

		select {
		case <-chResultExpired:
			bGameProceeds = false; //proceed to mmr calculation and bans for rq
		case oResult := <-chResult:
			MuGames.Lock();
			pGame.GameResult = oResult;
			MuGames.Unlock();
			//fmt.Printf("%v\n", oResult);
			if (oResult.GameEnded) {
				bGameProceeds = false;
			}
		}
	}



	//Game ended, settle results, destroy game
	MuGames.Lock();
	players.MuPlayers.Lock();

	pGame.State = StateGameEnded;
	oResult := pGame.GameResult;
	arFinalScores := rating.DetermineFinalScores(oResult, [2][]*players.EntPlayer{pGame.PlayersA, pGame.PlayersB});
	rating.UpdateMmr(oResult, arFinalScores, [2][]*players.EntPlayer{pGame.PlayersA, pGame.PlayersB});

	for _, pPlayer := range pGame.PlayersUnpaired {
		go database.UpdatePlayer(database.DatabasePlayer{
			SteamID64:			pPlayer.SteamID64,
			NicknameBase64:		pPlayer.NicknameBase64,
			Mmr:				pPlayer.Mmr,
			MmrUncertainty:		pPlayer.MmrUncertainty,
			Access:				pPlayer.Access,
			ProfValidated:		pPlayer.ProfValidated,
			RulesAccepted:		pPlayer.RulesAccepted,
			});
	}

	SetLastUpdated(pGame.PlayersUnpaired);
	Destroy(pGame);
	MuGames.Unlock();
	players.MuPlayers.Unlock();
}