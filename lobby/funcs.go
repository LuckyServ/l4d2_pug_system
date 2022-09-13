package lobby

import (
	//"fmt"
	"../utils"
	"../players"
	"../settings"
	"../games"
	"sort"
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
	//get list of online mmr's
	var arOnlineMmrs []int;
	var iOnlineCount int;
	games.MuGames.RLock();
	for _, pPlayer := range players.ArrayPlayers {
		if (pPlayer.ProfValidated && pPlayer.RulesAccepted && pPlayer.Access >= -1 && !pPlayer.IsIdle && (pPlayer.IsOnline || pPlayer.IsInLobby || IsFinishingGameSoon(pPlayer))) {
			arOnlineMmrs = append(arOnlineMmrs, pPlayer.Mmr);
		}
	}
	games.MuGames.RUnlock();
	if (pLobbyCreator.IsIdle) { //this is a very bad fix
		arOnlineMmrs = append(arOnlineMmrs, pLobbyCreator.Mmr);
	}
	iOnlineCount = len(arOnlineMmrs);

	//check if we already know the results
	if (iOnlineCount < settings.OnlineMmrRange * 2) {
		return -2000000000, 2000000000, nil;
	}

	//sort
	sort.Ints(arOnlineMmrs);

	//find index in the array of the lobby initial mmr
	iIndex := FindInitLobbyIndex(pLobbyCreator.Mmr, arOnlineMmrs);
	if (iIndex == -1) {
		return -1, -1, errors.New("Error calculating players range");
	}

	iOnlineMaxIdx := iOnlineCount - 1;
	iMaxMmrIdx := iIndex + settings.OnlineMmrRange;
	iMinMmrIdx := iIndex - settings.OnlineMmrRange;
	iMinMmr := -2000000000;
	iMaxMmr := 2000000000;
	if (iMaxMmrIdx > iOnlineMaxIdx) {
		iMinMmrIdx = iOnlineMaxIdx - (settings.OnlineMmrRange * 2);
		iMinMmr = arOnlineMmrs[iMinMmrIdx];
	} else if (iMinMmrIdx < settings.OnlineMmrRange - 1) {
		iMaxMmrIdx = settings.OnlineMmrRange * 2;
		iMaxMmr = arOnlineMmrs[iMaxMmrIdx];
	} else {
		iMinMmr = arOnlineMmrs[iMinMmrIdx];
		iMaxMmr = arOnlineMmrs[iMaxMmrIdx];
	}


	return iMinMmr, iMaxMmr, nil;
}

func FindInitLobbyIndex(iMmr int, arOnlineMmrs []int) int {
	iOnlineCount := len(arOnlineMmrs);
	if (iMmr < arOnlineMmrs[0]) { //shouldn't be possible, but just in case
		return 0;
	} else if (iMmr > arOnlineMmrs[iOnlineCount - 1]) { //shouldn't be possible, but just in case
		return iOnlineCount - 1;
	}

	iStartIdx := -1;
	iEndIdx := -1;
	for i, iOnlineMmr := range arOnlineMmrs {
		if (iMmr == iOnlineMmr) {
			iStartIdx = i;
			break;
		}
	}
	if (iStartIdx == -1) {
		return -1;
	}

	for i := iOnlineCount - 1; i >= iStartIdx; i-- {
		if (iMmr == arOnlineMmrs[i]) {
			iEndIdx = i;
			break;
		}
	}
	if (iEndIdx == -1 || iEndIdx < iStartIdx) {
		return -1;
	}

	return ((iStartIdx + iEndIdx) / 2);
}

func ChooseConfoglConfig(iMmr int) (settings.ConfoglConfig) {
	if (settings.BrokenMode) {
		return settings.ConfoglConfig{
			CodeName:		"default",
			Name:			"Default",
			MmrMax:			2000000000,
		};
	}
	iLen := len(settings.ArrayConfoglConfigsMmrs);
	if (iLen == 1) {
		return settings.MapConfoglConfigs[settings.ArrayConfoglConfigsMmrs[0]];
	}
	for i := 0; i < iLen; i++ {
		if (iMmr < settings.ArrayConfoglConfigsMmrs[i]) {
			return settings.MapConfoglConfigs[settings.ArrayConfoglConfigsMmrs[i]];
		}
	}
	return settings.ConfoglConfig{
		CodeName:		"zonemod",
		Name:			"nani?",
		MmrMax:			2000000000,
	}; //shouldn't happen
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