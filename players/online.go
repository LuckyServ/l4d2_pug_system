package players

import (
	"time"
)

func WatchOnline() {
	for {
		time.Sleep(5 * time.Second);

		i64CurTime := time.Now().UnixMilli();

		MuPlayers.Lock();
		for _, oPlayer := range ArrayPlayers {
			if (i64CurTime - oPlayer.LastActivity >= 60000/*1min*/) {
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
