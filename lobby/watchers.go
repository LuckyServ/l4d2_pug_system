package lobby

import (
	"../players"
	"../settings"
	"time"
)

func WatchLobbies() {
	for {
		time.Sleep(3 * time.Second);

		MuLobbies.Lock();
		players.MuPlayers.Lock();

		var arReadyLobbies []string;
		var arUnreadyPlayers, arTimedoutLobbiesPlayers, arJoinLobbyPlayers []*players.EntPlayer;
		
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

		for _, sLobbyID := range arReadyLobbies {
			_, bExists := MapLobbies[sLobbyID];
			if (bExists) {
				//Kill lobby and create game
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
		for _, pPlayer := range arJoinLobbyPlayers {
			JoinAny(pPlayer);
		}

		MuLobbies.Unlock();
		players.MuPlayers.Unlock();
	}
}