package rating

import (
	"../players"
	"../utils"
)


func Pair(arUnpairedPlayers []*players.EntPlayer) ([]*players.EntPlayer, []*players.EntPlayer) { //Games and Players must be locked outside

	//Sort
	bSorted := false;
	for !bSorted {
		bSorted = true;
		for i := 1; i < 8; i++ {
			if (arUnpairedPlayers[i].Mmr > arUnpairedPlayers[i - 1].Mmr) {
				arUnpairedPlayers[i], arUnpairedPlayers[i - 1] = arUnpairedPlayers[i - 1], arUnpairedPlayers[i]; //switch
				if (bSorted) {
					bSorted = false;
				}
			}
			if (!bSorted) {
				for i := 6; i >= 0; i-- {
					if (arUnpairedPlayers[i].Mmr < arUnpairedPlayers[i + 1].Mmr) {
						arUnpairedPlayers[i], arUnpairedPlayers[i + 1] = arUnpairedPlayers[i + 1], arUnpairedPlayers[i]; //switch
					}
				}
			}
		}
	}

	//Pair
	var arMostOfHappyDuos []int;
	var iCurMostOfHappyDuos int;
	for iBPicksTwo := 1; iBPicksTwo <= 5; iBPicksTwo += 2 {
		var arTeams = PlacePlayers(arUnpairedPlayers, iBPicksTwo);
		iHappyDuos := GetHappyDuoPlayers(arTeams);
		if (iHappyDuos > iCurMostOfHappyDuos) {
			iCurMostOfHappyDuos = iHappyDuos;
			arMostOfHappyDuos = make([]int, 0);
		}
		if (iHappyDuos == iCurMostOfHappyDuos) {
			arMostOfHappyDuos = append(arMostOfHappyDuos, iBPicksTwo);
		}
	}

	iRandInt, _ := utils.GetRandInt(0, len(arMostOfHappyDuos) - 1);
	var arTeams = PlacePlayers(arUnpairedPlayers, arMostOfHappyDuos[iRandInt]);

	return arTeams[0], arTeams[1];
}

func PlacePlayers(arUnpairedPlayers []*players.EntPlayer, iBPicksTwo int) ([2][]*players.EntPlayer) {
	var arTeams [2][]*players.EntPlayer;
	i := 0;
	iPicker := 0;
	for len(arTeams[0]) < 4 || len(arTeams[1]) < 4 {
		arTeams[iPicker] = append(arTeams[iPicker], arUnpairedPlayers[i]);
		if (iPicker == 0) {
			iPicker = 1;
		} else if (iBPicksTwo != i) {
			iPicker = 0;
		}
		i++;
	}
	return arTeams;
}

func GetHappyDuoPlayers(arTeams [2][]*players.EntPlayer) int {
	var iHappyDuoPlayers int;
	for iT := 0; iT < 2; iT++ {
		for iP := 0; iP < 4; iP++ {
			if (IsHappyDuoPlayer(arTeams[iT][iP], arTeams[iT])) {
				iHappyDuoPlayers++;
			}
		}
	}
	return iHappyDuoPlayers;
}

func IsHappyDuoPlayer(pPlayer *players.EntPlayer, arPlayers []*players.EntPlayer) bool {
	for _, pCheckPlayer := range arPlayers {
		if (pCheckPlayer.DuoWith == pPlayer.SteamID64) {
			return true;	
		}
	}
	return false;
}