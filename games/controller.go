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
	players.MuPlayers.Lock();
	MuGames.Lock();

	Create(pGame);
	pGame.State = StateCreated;
	SetLastUpdated(pGame.PlayersUnpaired);

	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Choose maps
	players.MuPlayers.Lock();
	MuGames.Lock();
	oCampaign := ChooseCampaign(pGame.PlayersUnpaired);
	pGame.CampaignName = oCampaign.Name;
	pGame.Maps = oCampaign.Maps;
	pGame.MapDownloadLink = oCampaign.DownloadLink;
	pGame.State = StateCampaignChosen;
	SetLastUpdated(pGame.PlayersUnpaired);
	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Pair players
	players.MuPlayers.Lock();
	MuGames.Lock();

	pGame.PlayersA, pGame.PlayersB = rating.Pair(pGame.PlayersUnpaired);

	pGame.State = StateTeamsPicked;
	SetLastUpdated(pGame.PlayersUnpaired);
	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Request pings
	players.MuPlayers.Lock();
	MuGames.Lock();
	for _, pPlayer := range pGame.PlayersUnpaired {
		pPlayer.GameServerPings = make(map[string]int);
	}
	pGame.State = StateWaitPings;
	SetLastUpdated(pGame.PlayersUnpaired);
	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Wait for pings
	time.Sleep(time.Duration(settings.MaxPingWait) * time.Second);


	//Cancel the ping request
	players.MuPlayers.Lock();
	MuGames.Lock();
	pGame.State = StateSelectServer;
	SetLastUpdated(pGame.PlayersUnpaired);
	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Select best available server based on pings and availability (a2s requests here)
	iTryCount := 0;
	for {

		arGameServers := GetGameServers(pGame.PlayersA, pGame.PlayersB); //Locks Players

		arGameServers = GetEmptyServers(arGameServers); //long execution, a2s queries here

		players.MuPlayers.Lock();
		MuGames.Lock();
		
		arGameServers = GetUnreservedServers(arGameServers);

		if (len(arGameServers) > 0) {
			pGame.ServerIP = arGameServers[0];
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
			if (iTryCount >= settings.AvailGameSrvsMaxTries) { //destroy game if too many tries
				players.MuPlayers.Lock();
				MuGames.Lock();

				go database.LogGame(database.DatabaseGameLog{
					ID:					pGame.ID,
					Valid:				false,
					CreatedAt:			pGame.CreatedAt,
					PlayersA:			Implode4Players(pGame.PlayersA),
					PlayersB:			Implode4Players(pGame.PlayersB),
					TeamAScores:		0,
					TeamBScores:		0,
					ConfoglConfig:		pGame.GameConfig.CodeName,
					CampaignName:		pGame.CampaignName,
					ServerIP:			"",
					Pings:				FormatPingsLog(pGame.PlayersUnpaired),
				});

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
		players.MuPlayers.Lock();
		MuGames.Lock();
		pGame.State = StateReadyUpExpired;
		SetLastUpdated(pGame.PlayersUnpaired);
		MuGames.Unlock();
		players.MuPlayers.Unlock();
	case <-chFullRUpReceived:
		players.MuPlayers.Lock();
		MuGames.Lock();
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
			players.MuPlayers.Lock();
			MuGames.Lock();

			go database.LogGame(database.DatabaseGameLog{
				ID:					pGame.ID,
				Valid:				false,
				CreatedAt:			pGame.CreatedAt,
				PlayersA:			Implode4Players(pGame.PlayersA),
				PlayersB:			Implode4Players(pGame.PlayersB),
				TeamAScores:		0,
				TeamBScores:		0,
				ConfoglConfig:		pGame.GameConfig.CodeName,
				CampaignName:		pGame.CampaignName,
				ServerIP:			pGame.ServerIP,
				Pings:				FormatPingsLog(pGame.PlayersUnpaired),
			});

			Destroy(pGame);
			MuGames.Unlock();
			players.MuPlayers.Unlock();
			return;
		case arReadyPlayers := <-chListOfReadyPlayers:
			if (len(arReadyPlayers) == 1 && (arReadyPlayers[0] == "none_ready" || arReadyPlayers[0] == "all_ready")) {
				//do nothing
			} else {
				//ban those who isnt ready
				if (len(arReadyPlayers) >= 4) { //at least 4 players must be ready for bans to happen
					var arBanReq []bans.EntAutoBanReq;
					players.MuPlayers.Lock();
					for _, pPlayer := range pGame.PlayersUnpaired {
						if (utils.GetStringIdxInArray(pPlayer.SteamID64, arReadyPlayers) == -1) {
							arBanReq = append(arBanReq, bans.EntAutoBanReq{
								SteamID64:			pPlayer.SteamID64,
								NicknameBase64:		pPlayer.NicknameBase64,
							});
							if (pGame.MapDownloadLink != "") {
								pPlayer.CustomMapsConfirmed = 0;
								go database.UpdatePlayer(database.DatabasePlayer{
									SteamID64:				pPlayer.SteamID64,
									NicknameBase64:			pPlayer.NicknameBase64,
									AvatarSmall:			pPlayer.AvatarSmall,
									AvatarBig:				pPlayer.AvatarBig,
									Mmr:					pPlayer.Mmr,
									MmrUncertainty:			pPlayer.MmrUncertainty,
									LastGameResult:			pPlayer.LastGameResult,
									Access:					pPlayer.Access,
									ProfValidated:			pPlayer.ProfValidated,
									RulesAccepted:			pPlayer.RulesAccepted,
									Twitch:					pPlayer.Twitch,
									CustomMapsConfirmed:	pPlayer.CustomMapsConfirmed,
									LastCampaignsPlayed:	strings.Join(pPlayer.LastCampaignsPlayed, "|"),
									});
								players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
							}
						}
					}
					players.MuPlayers.Unlock();
					time.Sleep(5 * time.Second);
					for _, oBanReq := range arBanReq {
						bans.ChanBanRQ <- oBanReq;
					}
				}
			}
			//Destroy the game
			players.MuPlayers.Lock();
			MuGames.Lock();

			go database.LogGame(database.DatabaseGameLog{
				ID:					pGame.ID,
				Valid:				false,
				CreatedAt:			pGame.CreatedAt,
				PlayersA:			Implode4Players(pGame.PlayersA),
				PlayersB:			Implode4Players(pGame.PlayersB),
				TeamAScores:		0,
				TeamBScores:		0,
				ConfoglConfig:		pGame.GameConfig.CodeName,
				CampaignName:		pGame.CampaignName,
				ServerIP:			pGame.ServerIP,
				Pings:				FormatPingsLog(pGame.PlayersUnpaired),
			});

			Destroy(pGame);
			MuGames.Unlock();
			players.MuPlayers.Unlock();
			return;
		case <-chFullRUpReceived:
			//Proceed to the next state
			players.MuPlayers.Lock();
			MuGames.Lock();
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
	players.MuPlayers.Lock();
	MuGames.Lock();

	pGame.State = StateGameEnded;
	oResult := pGame.GameResult;
	arFinalScores := rating.DetermineFinalScores(oResult, [2][]*players.EntPlayer{pGame.PlayersA, pGame.PlayersB});
	rating.UpdateMmr(oResult, arFinalScores, [2][]*players.EntPlayer{pGame.PlayersA, pGame.PlayersB});

	for _, pPlayer := range pGame.PlayersUnpaired {

		arCampaignsPlayed := pPlayer.LastCampaignsPlayed; //store played campaign
		iBuffer := utils.GetStringIdxInArray(pGame.CampaignName, arCampaignsPlayed);
		if (iBuffer != -1) {
			arCampaignsPlayed = append(arCampaignsPlayed[:iBuffer], arCampaignsPlayed[iBuffer+1:]...);
		}
		pPlayer.LastCampaignsPlayed = append([]string{pGame.CampaignName}, arCampaignsPlayed...);

		if (strings.HasPrefix(pPlayer.SteamID64, "7")) {
			go database.UpdatePlayer(database.DatabasePlayer{
				SteamID64:				pPlayer.SteamID64,
				NicknameBase64:			pPlayer.NicknameBase64,
				AvatarSmall:			pPlayer.AvatarSmall,
				AvatarBig:				pPlayer.AvatarBig,
				Mmr:					pPlayer.Mmr,
				MmrUncertainty:			pPlayer.MmrUncertainty,
				LastGameResult:			pPlayer.LastGameResult,
				Access:					pPlayer.Access,
				ProfValidated:			pPlayer.ProfValidated,
				RulesAccepted:			pPlayer.RulesAccepted,
				Twitch:					pPlayer.Twitch,
				CustomMapsConfirmed:	pPlayer.CustomMapsConfirmed,
				LastCampaignsPlayed:	strings.Join(pPlayer.LastCampaignsPlayed, "|"),
				});
		}
	}

	SetLastUpdated(pGame.PlayersUnpaired);
	players.I64LastPlayerlistUpdate = time.Now().UnixMilli();

	//Log game
	go database.LogGame(database.DatabaseGameLog{
		ID:					pGame.ID,
		Valid:				true,
		CreatedAt:			pGame.CreatedAt,
		PlayersA:			Implode4Players(pGame.PlayersA),
		PlayersB:			Implode4Players(pGame.PlayersB),
		TeamAScores:		arFinalScores[0],
		TeamBScores:		arFinalScores[1],
		ConfoglConfig:		pGame.GameConfig.CodeName,
		CampaignName:		pGame.CampaignName,
		ServerIP:			pGame.ServerIP,
		Pings:				FormatPingsLog(pGame.PlayersUnpaired),
	});

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
