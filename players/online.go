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
				I64LastPlayerlistUpdate = i64CurTime;
			}
		}
		MuPlayers.Unlock();

	}
}
