package rating

import (
	"../players"
	"time"
)


func Pair(arUnpairedPlayers []*players.EntPlayer) ([]*players.EntPlayer, []*players.EntPlayer) { //Games and Players must be locked outside

	//Sort
	arPlayersSorted := arUnpairedPlayers;
	bSorted := false;
	for !bSorted {
		bSorted = true;
		for i := 1; i < 8; i++ {
			if (arPlayersSorted[i].Mmr > arPlayersSorted[i - 1].Mmr) {
				arPlayersSorted[i], arPlayersSorted[i - 1] = arPlayersSorted[i - 1], arPlayersSorted[i]; //switch
				bSorted = false;
			}
			if (!bSorted) {
				for i := 6; i >= 0; i-- {
					if (arPlayersSorted[i].Mmr < arPlayersSorted[i + 1].Mmr) {
						arPlayersSorted[i], arPlayersSorted[i + 1] = arPlayersSorted[i + 1], arPlayersSorted[i]; //switch
					}
				}
			}
		}
	}

	//Pair
	var arTeams [2][]*players.EntPlayer;

	i := 0;
	iPicker := 0;
	iBPicksTwo := int(1 + ((time.Now().UnixNano() % int64(3)) * 2)); //random 1, 3, or 5
	for len(arTeams[0]) < 4 || len(arTeams[1]) < 4 {
		arTeams[iPicker] = append(arTeams[iPicker], arPlayersSorted[i]);
		if (iPicker == 0) {
			iPicker = 1;
		} else if (iBPicksTwo != i) {
			iPicker = 0;
		}
		i++;
	}

	return arTeams[0], arTeams[1];
}