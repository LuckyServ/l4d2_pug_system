package players

import (
	"time"
	"../settings"
)

func WatchOnline() {
	for {
		time.Sleep(5 * time.Second);

		i64CurTime := time.Now().UnixMilli();

		MuPlayers.Lock();
		for _, oPlayer := range ArrayPlayers {
			if (i64CurTime - oPlayer.LastActivity >= settings.OnlineTimeout) {
				if (oPlayer.IsOnline) {
					oPlayer.IsOnline = false;
					oPlayer.LastChanged = i64CurTime;
					I64LastPlayerlistUpdate = i64CurTime;
				}
			}
		}
		MuPlayers.Unlock();

	}
}
