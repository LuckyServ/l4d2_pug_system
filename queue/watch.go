package queue

import (
	"../players"
	"time"
	"../settings"
	"../games"
)

var i64MinWait int64 = 5 * 60 * 1000; //ms
var i64MaxWait int64 = 15 * 60 * 1000; //ms

func Watchers() {
	go WatchQueue();
}

func WatchQueue() {
	for {
		players.MuPlayers.Lock();

		if (BIsInReadyUp) {
			if (i64InReadyUpSince + settings.ReadyUpTimeout <= time.Now().UnixMilli()) {
				//kick unready
				KickUnready();

				//check if longest wait player is still there
				if (IPlayersCount >= 8 && pPlayerReadyUpReason == PLongestWaitPlayer) {
					SortQueueByWait();
					arTrimmedQueue := TrimQueue();
					SortTrimmedByMmr(arTrimmedQueue);

					//create games
					iNewGames := len(arTrimmedQueue) / 8;
					for iG := 0; iG < iNewGames; iG++ {

						pGame := &games.EntGame{
							ID:					<-games.ChanNewGameID,
							CreatedAt:			time.Now().UnixMilli(),
							State:				games.StateCreating,
						};

						var arGamePlayers []*players.EntPlayer;

						for iP := 0; iP < 8; iP++ {
							pPlayer := arTrimmedQueue[(8 * iG) + iP];
							arGamePlayers = append(arGamePlayers, pPlayer);
							Leave(pPlayer, true);
						}

						pGame.PlayersUnpaired = arGamePlayers;

						go games.Control(pGame);
					}
				}

				//stop readyup
				StopReadyUp();
				SetLastUpdated();
				players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
			}
		} else {
			if (IPlayersCount >= 8) {
				//check if minimal waiting time passed
				i64LongestWait := time.Now().UnixMilli() - PLongestWaitPlayer.InQueueSince;
				if (i64LongestWait >= i64MinWait) {
					if (i64LongestWait >= i64MaxWait) {
						RequestReadyUp();
					} else {
						//check if we have opportunity for more players to come soon
						if (games.IPlayersFinishingGameSoon > 0) {
						} else {
							RequestReadyUp();
						}
					}
				}
			}
		}

		players.MuPlayers.Unlock();
		time.Sleep(5 * time.Second);
	}
}