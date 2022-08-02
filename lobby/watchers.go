package lobby

import (
	"../players"
	"../settings"
	"../games"
	"time"
)

func WatchLobbies() {
	for {
		time.Sleep(3 * time.Second);

		MuLobbies.Lock();
		players.MuPlayers.Lock();

		var arReadyLobbies []string;
		var arUnreadyPlayers, arTimedoutLobbiesPlayers, arJoinLobbyPlayers, arGamePlayers []*players.EntPlayer;
		
		i64CurTime := time.Now().UnixMilli();

		for _, pLobby := range ArrayLobbies {
			if (pLobby.PlayerCount == 8 && pLobby.ReadyPlayers == 8) {
				arReadyLobbies = append(arReadyLobbies, pLobby.ID);
			} else if (pLobby.PlayerCount == 8 && pLobby.ReadyPlayers < 8 && i64CurTime - pLobby.ReadyUpSince >= settings.ReadyUpTimeout) {
				for _, pPlayer := range pLobby.Players {
					arUnreadyPlayers = append(arUnreadyPlayers, pPlayer);
				}
			} else if (pLobby.PlayerCount < 8 && i64CurTime - pLobby.CreatedAt >= settings.LobbyFillTimeout) {
				for _, pPlayer := range pLobby.Players {
					arTimedoutLobbiesPlayers = append(arTimedoutLobbiesPlayers, pPlayer);
				}
			}
		}


		for _, pPlayer := range arUnreadyPlayers {
			Leave(pPlayer);
			pPlayer.IsAutoSearching = false;
		}
		for _, pPlayer := range arTimedoutLobbiesPlayers {
			if (Leave(pPlayer) && pPlayer.IsAutoSearching) {
				arJoinLobbyPlayers = append(arJoinLobbyPlayers, pPlayer);
			}
		}
		//sort
		iSize := len(arJoinLobbyPlayers);
		if (iSize > 1) {
			bSorted := false;
			for !bSorted {
				bSorted = true;
				for i := 1; i < iSize; i++ {
					if (arJoinLobbyPlayers[i].AutoSearchingSince < arJoinLobbyPlayers[i - 1].AutoSearchingSince) {
						arJoinLobbyPlayers[i], arJoinLobbyPlayers[i - 1] = arJoinLobbyPlayers[i - 1], arJoinLobbyPlayers[i]; //switch
						bSorted = false;
					}
				}
			}
		}
		for _, pPlayer := range arJoinLobbyPlayers {
			JoinAny(pPlayer);
		}


		games.MuGames.Lock();
		for _, sLobbyID := range arReadyLobbies {
			pLobby, bExists := MapLobbies[sLobbyID];
			if (bExists) {

				pGame := &games.EntGame{
					ID:					<-games.ChanNewGameID,
					CreatedAt:			time.Now().UnixMilli(),
					GameConfig:			pLobby.GameConfig,
					State:				games.StateCreating,
					MmrMin:				pLobby.MmrMin,
					MmrMax:				pLobby.MmrMax,
				};

				for _, pPlayer := range pLobby.Players {
					arGamePlayers = append(arGamePlayers, pPlayer);
				}

				pGame.PlayersUnpaired = arGamePlayers;

				go games.Control(pGame);

			}
		}
		games.MuGames.Unlock();

		for _, pPlayer := range arGamePlayers {
			Leave(pPlayer);
			pPlayer.IsAutoSearching = false;
		}

		MuLobbies.Unlock();
		players.MuPlayers.Unlock();
	}
}