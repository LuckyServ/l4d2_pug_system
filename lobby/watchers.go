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

		var arReadyLobbies, arUnreadyPlayers []string;
		
		i64CurTime := time.Now().UnixMilli();

		for _, pLobby := range ArrayLobbies {
			if (pLobby.PlayerCount == 8 && pLobby.ReadyPlayers == 8) {
				arReadyLobbies = append(arReadyLobbies, pLobby.ID);
			} else if (pLobby.PlayerCount == 8 && pLobby.ReadyPlayers < 8 && i64CurTime - pLobby.ReadyUpSince >= settings.ReadyUpTimeout) {
				for _, pPlayer := range pLobby.Players {
					arUnreadyPlayers = append(arUnreadyPlayers, pPlayer.SteamID64);
				}
			}
		}

		for _, sLobbyID := range arReadyLobbies {
			_, bExists := MapLobbies[sLobbyID];
			if (bExists) {
				//Kill lobby and create game
			}
		}
		for _, sSteamID64 := range arUnreadyPlayers {
			pPlayer, bExists := players.MapPlayers[sSteamID64];
			if (bExists) {
				Leave(pPlayer);
			}
		}

		MuLobbies.Unlock();
		players.MuPlayers.Unlock();
	}
}