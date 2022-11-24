package queue

import (
	"../players"
	"time"
)

func SetLastUpdated() { //Players must be locked outside
	i64CurTime := time.Now().UnixMilli();
	for _, pPlayer := range arQueue {
		pPlayer.LastQueueChanged = i64CurTime;
	}
	players.I64LastPlayerlistUpdate = i64CurTime;
}

func FindPlayerInQueue(pPlayer *players.EntPlayer) int { //Players must be locked outside
	for i, pCheckPlayer := range arQueue {
		if (pCheckPlayer == pPlayer) {
			return i;	
		}
	}
	return -1;
}

func GetLongestWaitPlayer() (*players.EntPlayer) { //Players must be locked outside, len(arQueue) is guaranteed to be > 0
	var i64OldestJoin int64 = 9000000000000000000;
	var pOldestWaitPlayer *players.EntPlayer;
	for _, pPlayer := range arQueue {
		if (pPlayer.InQueueSince < i64OldestJoin) {
			i64OldestJoin = pPlayer.InQueueSince;
			pOldestWaitPlayer = pPlayer;
		}
	}
	return pOldestWaitPlayer;
}

func GetMmrDifferenceInQueue() int { //len guaranteed to be >=2
	var iMinMmr int = 2000000000;
	var iMaxMmr int = -2000000000;
	for _, pPlayer := range arQueue {
		if (pPlayer.Mmr < iMinMmr) {
			iMinMmr = pPlayer.Mmr;
		}
		if (pPlayer.Mmr > iMaxMmr) {
			iMaxMmr = pPlayer.Mmr;
		}
	}
	return iMaxMmr - iMinMmr;
}

func TrimQueue(arReadyOnly []*players.EntPlayer) ([]*players.EntPlayer) { //IPlayersCount must be >= 8 and arQueue sorted by wait time
	var arTrimmedQueue []*players.EntPlayer;
	iSize := len(arReadyOnly);
	iGamePlayers := iSize - (iSize % 8);
	for i := 0; i < iGamePlayers; i++ {
		arTrimmedQueue = append(arTrimmedQueue, arReadyOnly[i]);
	}
	return arTrimmedQueue;
}

func GetReadyPlayersOnly() ([]*players.EntPlayer) {
	var arReadyQueue []*players.EntPlayer;
	for _, pPlayer := range arQueue {
		if (pPlayer.IsReadyUpRequested && pPlayer.IsReadyConfirmed) {
			arReadyQueue = append(arReadyQueue, pPlayer);
		}
	}
	return arReadyQueue;
}

func SortTrimmedByMmr(arTrimmedQueue []*players.EntPlayer) {
	iSize := len(arTrimmedQueue);
	if (iSize > 1) {
		bSorted := false;
		for !bSorted {
			bSorted = true;
			for i := 1; i < iSize; i++ {
				if (arTrimmedQueue[i].Mmr < arTrimmedQueue[i - 1].Mmr) {
					arTrimmedQueue[i], arTrimmedQueue[i - 1] = arTrimmedQueue[i - 1], arTrimmedQueue[i]; //switch
					if (bSorted) {
						bSorted = false;
					}
				}
			}
			if (!bSorted) {
				for i := iSize - 2; i >= 0; i-- {
					if (arTrimmedQueue[i].Mmr > arTrimmedQueue[i + 1].Mmr) {
						arTrimmedQueue[i], arTrimmedQueue[i + 1] = arTrimmedQueue[i + 1], arTrimmedQueue[i]; //switch
					}
				}
			}
		}
	}
}

func SortQueueByWait(arReadyOnly []*players.EntPlayer) {
	iSize := len(arReadyOnly);
	if (iSize > 1) {
		bSorted := false;
		for !bSorted {
			bSorted = true;
			for i := 1; i < iSize; i++ {
				if (arReadyOnly[i].InQueueSince < arReadyOnly[i - 1].InQueueSince) {
					arReadyOnly[i], arReadyOnly[i - 1] = arReadyOnly[i - 1], arReadyOnly[i]; //switch
					if (bSorted) {
						bSorted = false;
					}
				}
			}
			if (!bSorted) {
				for i := iSize - 2; i >= 0; i-- {
					if (arReadyOnly[i].InQueueSince > arReadyOnly[i + 1].InQueueSince) {
						arReadyOnly[i], arReadyOnly[i + 1] = arReadyOnly[i + 1], arReadyOnly[i]; //switch
					}
				}
			}
		}
	}
}