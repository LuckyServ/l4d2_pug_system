package games

import (
	//"fmt"
	"../players"
	"../settings"
	"../rating"
	"../bans"
	"../utils"
	"time"
	"../database"
	"strings"
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
		pPlayer.GameServerPings = make(map[string]int, len(settings.GameServers));
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

		players.MuPlayers.Lock();
		sIPPORT := SelectBestAvailableServer(pGame.PlayersUnpaired); //unlocks players, but doesnt lock them

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

		players.MuPlayers.Lock();
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
				if (len(arReadyPlayers) >= 2) { //at least 2 players must be ready for bans to happen
					var arBanReq []bans.EntAutoBanReq;
					players.MuPlayers.RLock();
					for _, pPlayer := range pGame.PlayersUnpaired {
						if (utils.GetStringIdxInArray(pPlayer.SteamID64, arReadyPlayers) == -1) {
							arBanReq = append(arBanReq, bans.EntAutoBanReq{
								SteamID64:			pPlayer.SteamID64,
								NicknameBase64:		pPlayer.NicknameBase64,
							});
						}
					}
					players.MuPlayers.RUnlock();
					for _, oBanReq := range arBanReq {
						bans.ChanBanRQ <- oBanReq;
					}
				}
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
		if (strings.HasPrefix(pPlayer.SteamID64, "7")) {
			go database.UpdatePlayer(database.DatabasePlayer{
				SteamID64:			pPlayer.SteamID64,
				NicknameBase64:		pPlayer.NicknameBase64,
				Mmr:				pPlayer.Mmr,
				MmrUncertainty:		pPlayer.MmrUncertainty,
				LastGameResult:		pPlayer.LastGameResult,
				Access:				pPlayer.Access,
				ProfValidated:		pPlayer.ProfValidated,
				RulesAccepted:		pPlayer.RulesAccepted,
				});
		}
	}

	SetLastUpdated(pGame.PlayersUnpaired);
	players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
	Destroy(pGame);
	MuGames.Unlock();

	//store ban requests
	var arBanReq []bans.EntAutoBanReq;
	for _, sSteamID64 := range oResult.AbsentPlayers {
		pPlayer, bFound := players.MapPlayers[sSteamID64];
		if (bFound) {
			arBanReq = append(arBanReq, bans.EntAutoBanReq{
				SteamID64:			sSteamID64,
				NicknameBase64:		pPlayer.NicknameBase64,
			});
		}
	}

	players.MuPlayers.Unlock();

	time.Sleep(1 * time.Second);

	//Ban ragequitters
	for _, oBanReq := range arBanReq {
		bans.ChanBanRQ <- oBanReq;
	}
}
