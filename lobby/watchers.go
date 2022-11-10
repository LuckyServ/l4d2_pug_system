package lobby

import (
	"../players"
	"../settings"
	"../games"
	"time"
)

var ChanUniqueTime = make(chan int64);

func Watchers() {
	go WatchLobbies();
	go SortLobbies();
	go WatchUniqueTimeChan();
	go RecreateLobbies();
}

func WatchUniqueTimeChan() {
	for {
		select {
		case ChanUniqueTime <- func()(int64) {
			return time.Now().UnixMilli();
		}():
		}
		time.Sleep(1 * time.Millisecond);
	}
}

func SortLobbies() {
	for {
		time.Sleep(3 * time.Second);

		bEdited := false;
		MuLobbies.Lock();
		iSize := len(ArrayLobbies);
		if (iSize > 1) {
			bSorted := false;
			for !bSorted {
				bSorted = true;
				for i := 1; i < iSize; i++ {
					if (ArrayLobbies[i].CreatedAt < ArrayLobbies[i - 1].CreatedAt) {
						ArrayLobbies[i], ArrayLobbies[i - 1] = ArrayLobbies[i - 1], ArrayLobbies[i]; //switch
						if (!bEdited) {
							bEdited = true;
						}
						if (bSorted) {
							bSorted = false;
						}
					}
				}
				if (!bSorted) {
					for i := iSize - 2; i >= 0; i-- {
						if (ArrayLobbies[i].CreatedAt > ArrayLobbies[i + 1].CreatedAt) {
							ArrayLobbies[i], ArrayLobbies[i + 1] = ArrayLobbies[i + 1], ArrayLobbies[i]; //switch
							if (!bEdited) {
								bEdited = true;
							}
						}
					}
				}
			}
		}

		MuLobbies.Unlock();
		if (bEdited) {
			I64LastLobbyListUpdate = time.Now().UnixMilli();
		}
	}
}

func RecreateLobbies() {
	for {
		time.Sleep(30 * time.Second);

		players.MuPlayers.Lock();
		MuLobbies.Lock();


		var arRecreateLobbiesPlayers, arJoinLobbyPlayers []*players.EntPlayer;

		for _, pLobby := range ArrayLobbies {
			if (pLobby.PlayerCount < 8) {
				for _, pPlayer := range pLobby.Players {
					arRecreateLobbiesPlayers = append(arRecreateLobbiesPlayers, pPlayer);
				}
			}
		}

		for _, pPlayer := range arRecreateLobbiesPlayers {
			if (Leave(pPlayer, true) && pPlayer.IsAutoSearching) {
				arJoinLobbyPlayers = append(arJoinLobbyPlayers, pPlayer);
			}
		}

		iSize := len(arJoinLobbyPlayers);
		if (iSize > 1) {
			bSorted := false;
			for !bSorted {
				bSorted = true;
				for i := 1; i < iSize; i++ {
					if (arJoinLobbyPlayers[i].AutoSearchingSince < arJoinLobbyPlayers[i - 1].AutoSearchingSince) {
						arJoinLobbyPlayers[i], arJoinLobbyPlayers[i - 1] = arJoinLobbyPlayers[i - 1], arJoinLobbyPlayers[i]; //switch
						if (bSorted) {
							bSorted = false;
						}
					}
				}
				if (!bSorted) {
					for i := iSize - 2; i >= 0; i-- {
						if (arJoinLobbyPlayers[i].AutoSearchingSince > arJoinLobbyPlayers[i + 1].AutoSearchingSince) {
							arJoinLobbyPlayers[i], arJoinLobbyPlayers[i + 1] = arJoinLobbyPlayers[i + 1], arJoinLobbyPlayers[i]; //switch
						}
					}
				}
			}
		}
		for _, pPlayer := range arJoinLobbyPlayers {
			JoinAny(pPlayer);
		}


		players.MuPlayers.Unlock();
		MuLobbies.Unlock();
	}
}

func WatchLobbies() {
	for {
		time.Sleep(3 * time.Second);

		players.MuPlayers.Lock();
		MuLobbies.Lock();

		var arReadyLobbies []string;
		var arUnreadyPlayers, arOfflinePlayers, arGamePlayers []*players.EntPlayer;
		
		i64CurTime := time.Now().UnixMilli();

		for _, pLobby := range ArrayLobbies {
			if (pLobby.PlayerCount == 8 && pLobby.ReadyPlayers == 8) {
				arReadyLobbies = append(arReadyLobbies, pLobby.ID);
			} else if (pLobby.PlayerCount == 8 && pLobby.ReadyPlayers < 8 && i64CurTime - pLobby.ReadyUpSince >= settings.ReadyUpTimeout) {
				for _, pPlayer := range pLobby.Players {
					if (!pPlayer.IsReadyInLobby) {
						arUnreadyPlayers = append(arUnreadyPlayers, pPlayer);
					}
				}
			} else if (pLobby.PlayerCount < 8) {
				for _, pPlayer := range pLobby.Players {
					if (!pPlayer.IsOnline) {
						arOfflinePlayers = append(arOfflinePlayers, pPlayer);
					}
				}
			}
		}

		for _, pPlayer := range arOfflinePlayers {
			Leave(pPlayer, false);
		}
		for _, pPlayer := range arUnreadyPlayers {
			Leave(pPlayer, false);
		}


		if (len(arReadyLobbies) > 0) {
			games.MuGames.Lock();
			for _, sLobbyID := range arReadyLobbies {
				pLobby, bExists := MapLobbies[sLobbyID];
				if (bExists) {

					pGame := &games.EntGame{
						ID:					<-games.ChanNewGameID,
						CreatedAt:			time.Now().UnixMilli(),
						State:				games.StateCreating,
					};

					for _, pPlayer := range pLobby.Players {
						arGamePlayers = append(arGamePlayers, pPlayer);
					}

					pGame.PlayersUnpaired = arGamePlayers;

					go games.Control(pGame);

				}
			}
			games.MuGames.Unlock();
		}

		for _, pPlayer := range arGamePlayers {
			Leave(pPlayer, true);
			pPlayer.IsAutoSearching = false;
		}

		players.MuPlayers.Unlock();
		MuLobbies.Unlock();
	}
}