package queue

import (
	"../players"
	"time"
)

var arQueue []*players.EntPlayer;
var NewGamesBlocked bool;
var IPlayersCount int;
var bIsInReadyUp bool;
var i64InReadyUpSince int64;
var PLongestWaitPlayer *players.EntPlayer;
var pPlayerReadyUpReason *players.EntPlayer;
var IReadyPlayers int;

var i64CooldownForReadyUpLeave int64 = 5 * 60 * 1000; //ms
var i64CooldownForLeave int64 = 5 * 1000; //ms


func Join(pPlayer *players.EntPlayer) { //Players must be locked outside

	if (pPlayer.IsInQueue) { //repeat critical check
		return;
	}

	arQueue = append(arQueue, pPlayer);
	IPlayersCount++;
	pPlayer.IsInQueue = true;
	pPlayer.InQueueSince = time.Now().UnixMilli();
	pPlayer.IsReadyUpRequested = false;
	pPlayer.IsReadyConfirmed = false;
	if (IPlayersCount == 1) {
		PLongestWaitPlayer = pPlayer;
	}
	SetLastUpdated();
}

func Leave(pPlayer *players.EntPlayer, bGameStart bool) { //Players must be locked outside
	iPlayer := FindPlayerInQueue(pPlayer);
	if (iPlayer != -1) {
		arQueue[iPlayer] = arQueue[len(arQueue)-1];
		arQueue = arQueue[:len(arQueue)-1];
		IPlayersCount--;

		if (bGameStart) {
			pPlayer.NextQueueingAllowed = 0;
		} else if (pPlayer.IsReadyUpRequested) {
			pPlayer.NextQueueingAllowed = time.Now().UnixMilli() + i64CooldownForReadyUpLeave;
		} else {
			pPlayer.NextQueueingAllowed = time.Now().UnixMilli() + i64CooldownForLeave;
		}

		pPlayer.IsInQueue = false;
		pPlayer.InQueueSince = 0;
		pPlayer.IsReadyUpRequested = false;
		pPlayer.IsReadyConfirmed = false;

		if (pPlayer == PLongestWaitPlayer) {
			if (IPlayersCount == 0) {
				PLongestWaitPlayer = nil;
			} else {
				PLongestWaitPlayer = GetLongestWaitPlayer();
			}
		}

		SetLastUpdated();
	}
}

func ReadyUp(pPlayer *players.EntPlayer) { //Players must be locked outside
	pPlayer.IsReadyConfirmed = true;
	SetLastUpdated();
}

func RequestReadyUp() { //queue is >= 8 players
	bIsInReadyUp = true;
	i64InReadyUpSince = time.Now().UnixMilli();
	pPlayerReadyUpReason = PLongestWaitPlayer;
	IReadyPlayers = 0;
	for _, pPlayer := range arQueue {
		pPlayer.IsReadyUpRequested = true;
		pPlayer.IsReadyConfirmed = false;
	}
	SetLastUpdated();
}

func StopReadyUp() {
	bIsInReadyUp = false;
	i64InReadyUpSince = 0;
	IReadyPlayers = 0;
	pPlayerReadyUpReason = nil;
	for _, pPlayer := range arQueue {
		pPlayer.IsReadyUpRequested = false;
		pPlayer.IsReadyConfirmed = false;
	}
}

func KickUnready() {
	var arKickPlayers []*players.EntPlayer;
	for _, pPlayer := range arQueue {
		if (pPlayer.IsReadyUpRequested && !pPlayer.IsReadyConfirmed) {
			arKickPlayers = append(arKickPlayers, pPlayer);
		}
	}
	for _, pPlayer := range arKickPlayers {
		Leave(pPlayer, false);
	}
}