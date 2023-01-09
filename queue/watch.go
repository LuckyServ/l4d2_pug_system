package queue

import (
	"../players"
	"time"
	"../settings"
	"../games"
	"math/big"
	"strings"
	"fmt"
)

var I64MaxQueueWait int64 = 15 * 60 * 1000; //ms
var GenInviteCode = make(chan string);

func Watchers() {
	go WatchQueue();
	go WatchKickOffline();
	go DuoOfferIDGenerator();
}

func WatchKickOffline() {
	for {
		players.MuPlayers.Lock();
		KickOffline();
		players.MuPlayers.Unlock();
		time.Sleep(5 * time.Second);
	}
}

func WatchQueue() {
	for {
		players.MuPlayers.Lock();

		if (BIsInReadyUp) {
			if (i64InReadyUpSince + settings.ReadyUpTimeout <= time.Now().UnixMilli()) {
				//kick unready
				KickUnready();

				arReadyOnly := GetReadyPlayersOnly();

				//check if longest wait player is still there
				if (len(arReadyOnly) >= 8 && pPlayerReadyUpReason == PLongestWaitPlayer) {
					arTrimmedQueue := TrimQueue(arReadyOnly);
					if (len(arTrimmedQueue) == 0) {
						bWaitingForSinglePlayer = true;
					} else {
						arTrimmedQueue = SortTrimmedByMmr(arTrimmedQueue);

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
				}

				//stop readyup
				StopReadyUp();
				SetLastUpdated();
			}
		} else {
			if (IPlayersCount >= 8 && !bWaitingForSinglePlayer) {
				//check if minimal waiting time passed
				i64LongestWait := time.Now().UnixMilli() - PLongestWaitPlayer.InQueueSince;
				if (i64LongestWait >= I64MaxQueueWait) {
					RequestReadyUp();
				}
			}
		}

		players.MuPlayers.Unlock();
		time.Sleep(5 * time.Second);
	}
}

func DuoOfferIDGenerator() {
	for {
		select {
		case GenInviteCode <- func()(string) {
			return strings.ToUpper(fmt.Sprintf("L4D2C%s", big.NewInt(time.Now().UnixNano()).Text(36)));
		}():
		}
		time.Sleep(1 * time.Nanosecond);
	}
}