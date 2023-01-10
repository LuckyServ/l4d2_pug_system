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

func FindPlayerInArray(pPlayer *players.EntPlayer, arPlayers []*players.EntPlayer) int { //Players must be locked outside
	for i, pCheckPlayer := range arPlayers {
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

func TrimQueue(arReadyOnly []*players.EntPlayer) ([]*players.EntPlayer) { //IPlayersCount must be >= 8 and arQueue sorted by wait time
	var arTrimmedQueue []*players.EntPlayer;
	iSize := len(arReadyOnly);
	iGamePlayers := iSize - (iSize % 8);

	if (iSize > iGamePlayers && AreDuoQueued(arReadyOnly[iGamePlayers - 1], arReadyOnly[iGamePlayers])) {
		//cant trim, players are duo queued, need to move them
		iNewerSinglePlayer := GetNewerSinglePlayer(arReadyOnly, iGamePlayers);
		if (iNewerSinglePlayer == -1) {
			//cant move, return failure
			return arTrimmedQueue;
		} else {
			//move single player in queue
			pNewerSinglePlayer := arReadyOnly[iNewerSinglePlayer];
			arReadyOnly = append(arReadyOnly[:iNewerSinglePlayer], arReadyOnly[iNewerSinglePlayer+1:]...);
			arReadyOnly = append(arReadyOnly[:iGamePlayers-1], append([]*players.EntPlayer{pNewerSinglePlayer}, arReadyOnly[iGamePlayers-1:]...)...);
		}
	}

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

func SortTrimmedByMmr(arTrimmedQueue []*players.EntPlayer) []*players.EntPlayer {

	//convert players into arrays of 1 or 2 players
	var arOfArraysOfPlayers [][]*players.EntPlayer;
	iSize := len(arTrimmedQueue);
	for i := 0; i < iSize - 1; i++ {
		if (AreDuoQueued(arTrimmedQueue[i], arTrimmedQueue[i + 1])) {
			arOfArraysOfPlayers = append(arOfArraysOfPlayers, []*players.EntPlayer{arTrimmedQueue[i], arTrimmedQueue[i + 1]});
			i++;
		} else {
			arOfArraysOfPlayers = append(arOfArraysOfPlayers, []*players.EntPlayer{arTrimmedQueue[i]});
		}
	}
	if (arTrimmedQueue[iSize - 1].DuoWith == "") {
		arOfArraysOfPlayers = append(arOfArraysOfPlayers, []*players.EntPlayer{arTrimmedQueue[iSize - 1]});
	}

	//sort by ranked/unranked and mmr
	iArSize := len(arOfArraysOfPlayers);
	if (iArSize > 1) {
		bSorted := false;
		for !bSorted {
			bSorted = true;
			for i := 1; i < iArSize; i++ {
				if (GetAvgMmr(arOfArraysOfPlayers[i]) < GetAvgMmr(arOfArraysOfPlayers[i - 1])) {
					arOfArraysOfPlayers[i], arOfArraysOfPlayers[i - 1] = arOfArraysOfPlayers[i - 1], arOfArraysOfPlayers[i]; //switch
					if (bSorted) {
						bSorted = false;
					}
				}
			}
			if (!bSorted) {
				for i := iArSize - 2; i >= 0; i-- {
					if (GetAvgMmr(arOfArraysOfPlayers[i]) > GetAvgMmr(arOfArraysOfPlayers[i + 1])) {
						arOfArraysOfPlayers[i], arOfArraysOfPlayers[i + 1] = arOfArraysOfPlayers[i + 1], arOfArraysOfPlayers[i]; //switch
					}
				}
			}
		}
		bSorted = false;
		for !bSorted {
			bSorted = true;
			for i := 1; i < iArSize; i++ {
				if (!IsGroupRanked(arOfArraysOfPlayers[i]) && IsGroupRanked(arOfArraysOfPlayers[i - 1])) {
					arOfArraysOfPlayers[i], arOfArraysOfPlayers[i - 1] = arOfArraysOfPlayers[i - 1], arOfArraysOfPlayers[i]; //switch
					if (bSorted) {
						bSorted = false;
					}
				}
			}
			if (!bSorted) {
				for i := iArSize - 2; i >= 0; i-- {
					if (IsGroupRanked(arOfArraysOfPlayers[i]) && !IsGroupRanked(arOfArraysOfPlayers[i + 1])) {
						arOfArraysOfPlayers[i], arOfArraysOfPlayers[i + 1] = arOfArraysOfPlayers[i + 1], arOfArraysOfPlayers[i]; //switch
					}
				}
			}
		}
	}

	//convert back into single-dimensional array
	arTrimmedQueue = make([]*players.EntPlayer, 0);
	for _, arOfPlayers := range arOfArraysOfPlayers {
		for _, pPlayer := range arOfPlayers {
			arTrimmedQueue = append(arTrimmedQueue, pPlayer);
		}
	}

	//check if there are duo conflicts found, fix them
	iNewGames := iSize / 8;
	for iG := 0; iG < iNewGames - 1; iG++ {
		pPlayer1 := arTrimmedQueue[(8 * iG) + 7];
		pPlayer2 := arTrimmedQueue[(8 * (iG + 1))];
		if (AreDuoQueued(pPlayer1, pPlayer2)) {
			iNearestByMmrSinglePlayer := GetNearestByMmrSinglePlayer(arTrimmedQueue, (8 * iG) + 7);
			iNewPlace := func()(int) {
				if (iNearestByMmrSinglePlayer > (8 * (iG + 1))) {
					return (8 * iG) + 7;
				}
				return (8 * (iG + 1));
			}();
			//move single player in queue
			pNewerSinglePlayer := arTrimmedQueue[iNearestByMmrSinglePlayer];
			arTrimmedQueue = append(arTrimmedQueue[:iNearestByMmrSinglePlayer], arTrimmedQueue[iNearestByMmrSinglePlayer+1:]...);
			arTrimmedQueue = append(arTrimmedQueue[:iNewPlace], append([]*players.EntPlayer{pNewerSinglePlayer}, arTrimmedQueue[iNewPlace:]...)...);
		}
	}
	return arTrimmedQueue;
}

func AreDuoQueued(pPlayer1 *players.EntPlayer, pPlayer2 *players.EntPlayer) bool {
	if (pPlayer1.SteamID64 == pPlayer2.DuoWith) {
		return true;
	}
	return false;
}

func GetNearestByMmrSinglePlayer(arTrimmedQueue []*players.EntPlayer, iFirstInDuo int) int {
	iSize := len(arTrimmedQueue);
	
	var iCurrentNearestLowerMmrPlayer int = -1
	var iCurrentNearestHigherMmrPlayer int = -1;
	for i := iFirstInDuo - 1; i >= 0; i-- {
		if (arTrimmedQueue[i].DuoWith == "") {
			iCurrentNearestLowerMmrPlayer = i;
			break;
		}
	}
	for i := iFirstInDuo + 2; i < iSize; i++ {
		if (arTrimmedQueue[i].DuoWith == "") {
			iCurrentNearestHigherMmrPlayer = i;
			break;
		}
	}

	iDuoMmr := GetAvgMmr([]*players.EntPlayer{arTrimmedQueue[iFirstInDuo], arTrimmedQueue[iFirstInDuo + 1]});
	if (iCurrentNearestLowerMmrPlayer == -1) {
		return iCurrentNearestHigherMmrPlayer;
	} else if (iCurrentNearestHigherMmrPlayer == -1) {
		return iCurrentNearestLowerMmrPlayer;
	} else if (iDuoMmr - arTrimmedQueue[iCurrentNearestLowerMmrPlayer].Mmr < arTrimmedQueue[iCurrentNearestHigherMmrPlayer].Mmr - iDuoMmr) {
		return iCurrentNearestLowerMmrPlayer;
	} else {
		return iCurrentNearestHigherMmrPlayer;
	}

	return -1;
}

func GetNewerSinglePlayer(arReadyOnly []*players.EntPlayer, iGamePlayers int) int {
	iSize := len(arReadyOnly);
	for i := iGamePlayers + 1; i < iSize; i++ {
		if (arReadyOnly[i].DuoWith == "") {
			return i;
		}
	}
	return -1;
}

func GetAvgMmr(arPlayers []*players.EntPlayer) int {
	iSize := len(arPlayers);
	var iMmrSum int;
	if (iSize > 0) {
		for _, pPlayer := range arPlayers {
			iMmrSum = iMmrSum + pPlayer.Mmr;
		}
		return (iMmrSum / iSize);
	}
	return 0;
}

func IsGroupRanked(arPlayers []*players.EntPlayer) bool {
	iSize := len(arPlayers);
	if (iSize > 0) {
		for _, pPlayer := range arPlayers {
			if (players.GetMmrGrade(pPlayer) == 0) {
				return false;
			}
		}
		return true;
	}
	return false;
}
