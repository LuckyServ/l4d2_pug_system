package players

import (
	"time"
	"../settings"
)

func WatchOnline() {
	for {
		time.Sleep(5 * time.Second);

		MuPlayers.Lock();
		i64CurTime := time.Now().UnixMilli();
		for _, pPlayer := range ArrayPlayers {
			if (pPlayer.IsOnline && i64CurTime - pPlayer.LastActivity >= settings.OnlineTimeout) { //offline
				pPlayer.IsOnline = false;
				pPlayer.IsIdle = false;
				I64LastPlayerlistUpdate = i64CurTime;
			} else if (pPlayer.IsOnline && !pPlayer.IsIdle && !pPlayer.IsInGame && !pPlayer.IsInLobby &&
				(i64CurTime - pPlayer.OnlineSince) >= settings.IdleTimeout &&
				(i64CurTime - pPlayer.LastLobbyActivity) >= settings.IdleTimeout &&
				(i64CurTime - pPlayer.LastGameActivity) >= settings.IdleTimeout) { //idle
				pPlayer.IsIdle = true;
				I64LastPlayerlistUpdate = i64CurTime;
			}
		}
		MuPlayers.Unlock();

	}
}
