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
	var arTeams [2][]*players.EntPlayer;

	i := 0;
	iPicker := 0;
	iRandInt, _ := utils.GetRandInt(0, 2);
	iBPicksTwo := (iRandInt * 2) + 1; //random 1, 3, or 5
	for len(arTeams[0]) < 4 || len(arTeams[1]) < 4 {
		arTeams[iPicker] = append(arTeams[iPicker], arUnpairedPlayers[i]);
		if (iPicker == 0) {
			iPicker = 1;
		} else if (iBPicksTwo != i) {
			iPicker = 0;
		}
		i++;
	}

	return arTeams[0], arTeams[1];
}