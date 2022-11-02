package lobby

import (
	//"fmt"
	"../utils"
	"../players"
	"../settings"
	"../games"
	"errors"
)


func GenerateID() string { //MuLobbies must be blocked outside
	var sLobbyID string;
	var iLength int = 2;
	bIsUnique := false;
	for !bIsUnique {
		sLobbyID, _ = utils.GenerateRandomString(iLength, "123456789");
		bIsUnique = true;
		for _, oLobby := range ArrayLobbies {
			if (sLobbyID == oLobby.ID) {
				bIsUnique = false;
				iLength++;
				break;
			}
		}
	}
	return sLobbyID;
}

func CalcMmrLimits(pLobbyCreator *players.EntPlayer) (int, int, error) { //MuPlayers must be locked outside

	players.SortPlayers();

	iPlayer := FindPlayerIndex(pLobbyCreator);
	if (iPlayer == -1) {
		return -1, -1, errors.New("Error calculating players range");
	}

	iCount := len(players.ArrayPlayers);
	if (iCount < (settings.OnlineMmrRange * 2) + 1) {
		return -2000000000, 2000000000, nil;
	}

	var iMinMmr int = -2000000000;
	var iMaxMmr int = 2000000000;
	var iCurRangeMin, iCurRangeMax int;
	var iEndedOnPlayerMin, iEndedOnPlayerMax int;

	games.MuGames.RLock();

	if (iPlayer >= 1) {
		for i := iPlayer - 1; i >= 0; i-- {
			if (PlayerSuitableForMmrRangeCalc(players.ArrayPlayers[i])) {
				iMaxMmr = players.ArrayPlayers[i].Mmr;
				iEndedOnPlayerMax = i;
				iCurRangeMax++;
			}
			if (iCurRangeMax >= settings.OnlineMmrRange) {
				break;
			}
		}
	}

	if (iPlayer < (iCount - 1)) {
		for i := iPlayer + 1; i < iCount; i++ {
			if (PlayerSuitableForMmrRangeCalc(players.ArrayPlayers[i])) {
				iMinMmr = players.ArrayPlayers[i].Mmr;
				iEndedOnPlayerMin = i;
				iCurRangeMin++;
			}
			if (iCurRangeMin >= settings.OnlineMmrRange) {
				break;
			}
		}
	}

	var iMoreMaxMmrNeeded, iMoreMinMmrNeeded int;
	if (iCurRangeMin < settings.OnlineMmrRange) {
		iMoreMaxMmrNeeded = settings.OnlineMmrRange - iCurRangeMin;
		iMinMmr = -2000000000;
	}
	if (iCurRangeMax < settings.OnlineMmrRange) {
		iMoreMinMmrNeeded = settings.OnlineMmrRange - iCurRangeMax;
		iMaxMmr = 2000000000;
	}

	var iCurRangeMinAdd, iCurRangeMaxAdd int;

	if (iMoreMaxMmrNeeded > 0 && iEndedOnPlayerMax - iMoreMaxMmrNeeded >= 0) {
		for i := iEndedOnPlayerMax - 1; i >= 0; i-- {
			if (PlayerSuitableForMmrRangeCalc(players.ArrayPlayers[i])) {
				iMaxMmr = players.ArrayPlayers[i].Mmr;
				iCurRangeMaxAdd++;
			}
			if (iCurRangeMaxAdd >= iMoreMaxMmrNeeded) {
				break;
			}
		}
	}
	if (iMoreMinMmrNeeded > 0 && iEndedOnPlayerMin + iMoreMinMmrNeeded < iCount) {
		for i := iEndedOnPlayerMin + 1; i < iCount; i++ {
			if (PlayerSuitableForMmrRangeCalc(players.ArrayPlayers[i])) {
				iMinMmr = players.ArrayPlayers[i].Mmr;
				iCurRangeMinAdd++;
			}
			if (iCurRangeMinAdd >= iMoreMinMmrNeeded) {
				break;
			}
		}
	}

	games.MuGames.RUnlock();

	
	return iMinMmr, iMaxMmr, nil;
}

func PlayerSuitableForMmrRangeCalc(pPlayer *players.EntPlayer) bool { //MuPlayers and MyGames must be locked outside
	if (!pPlayer.ProfValidated || !pPlayer.RulesAccepted || pPlayer.Access < 0) {
		return false;
	}
	if ((!pPlayer.IsOnline || pPlayer.IsIdle) && !pPlayer.IsInLobby && !pPlayer.IsInGame) {
		return false;
	}
	if (pPlayer.IsInGame && !IsFinishingGameSoon(pPlayer)) {
		return false;
	}
	return true;
}

func FindPlayerIndex(pSearchedPlayer *players.EntPlayer) int { //MuPlayers must be locked outside
	for i, pPlayer := range players.ArrayPlayers {
		if (pSearchedPlayer == pPlayer) {
			return i;
		}
	}
	return -1;
}

func GetJoinableLobbies(iMmr int) []*EntLobby { //MuLobbies must be locked outside
	var arLobbies []*EntLobby;
	for _, pLobby := range ArrayLobbies {
		if (iMmr >= pLobby.MmrMin && iMmr <= pLobby.MmrMax && pLobby.PlayerCount < 8/*hardcoded for 4v4*/) {
			arLobbies = append(arLobbies, pLobby);
		}
	}
	return arLobbies;
}

func IsFinishingGameSoon(pPlayer *players.EntPlayer) bool { //Players and Games must be locked outside
	if (!pPlayer.IsInGame) {
		return false;
	}
	oGameResult := games.MapGames[pPlayer.GameID].GameResult;
	if (oGameResult.IsLastMap && oGameResult.CurrentHalf == 2) {
		return true;
	}
	return false;
}