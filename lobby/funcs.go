package lobby

import (
	//"fmt"
	"../utils"
	"../players"
	"../settings"
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

func CalcMmrLimits(iMmr int) (int, int, error) { //MuPlayers must be locked outside
	//get list of online mmr's
	var arOnlineMmrs []int;
	var iOnlineCount int;
	for _, oPlayer := range players.ArrayPlayers {
		if ((oPlayer.IsOnline || oPlayer.IsInLobby) && !oPlayer.IsInGame && oPlayer.ProfValidated && oPlayer.RulesAccepted && oPlayer.Access >= -1/*not banned*/) {
			arOnlineMmrs = append(arOnlineMmrs, oPlayer.Mmr);
		}
	}
	iOnlineCount = len(arOnlineMmrs);

	//check if we already know the results
	if (iOnlineCount < settings.OnlineMmrRange * 2) {
		return -2000000000, 2000000000, nil;
	}

	//sort
	sort.Ints(arOnlineMmrs);

	//find index in the array of the lobby initial mmr
	iIndex := FindInitLobbyIndex(iMmr, arOnlineMmrs);
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

func ChooseConfoglConfig(iMmr int) string {
	iLen := len(settings.ArrayConfoglConfigsMmrs);
	if (iLen == 1) {
		return settings.MapConfoglConfigs[settings.ArrayConfoglConfigsMmrs[0]];
	}
	for i := 0; i < iLen; i++ {
		if (iMmr < settings.ArrayConfoglConfigsMmrs[i]) {
			return settings.MapConfoglConfigs[settings.ArrayConfoglConfigsMmrs[i]];
		}
	}
	return "zonemod"; //shouldn't happen
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
