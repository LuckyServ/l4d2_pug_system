package rating

import (
	"../players"
	"../utils"
	"../settings"
)


func Pair(arUnpairedPlayers []*players.EntPlayer) ([]*players.EntPlayer, []*players.EntPlayer) { //Games and Players must be locked outside

	//Sort players
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

	//Replace integers with Players in variants
	var arVariantsP [][2][]*players.EntPlayer;
	for _, arVariantI := range arPairVariants {
		var arVariantP [2][]*players.EntPlayer;
		for iT := 0; iT < 2; iT++ {
			for iP := 0; iP < 4; iP++ {
				arVariantP[iT] = append(arVariantP[iT], arUnpairedPlayers[arVariantI[iT][iP]]);
			}
		}
		arVariantsP = append(arVariantsP, arVariantP);
	}

	//Sort variants by mmr difference
	iSize := len(arVariantsP);
	if (iSize > 1) {
		bSorted := false;
		for !bSorted {
			bSorted = true;
			for i := 1; i < iSize; i++ {
				if (GetMmrDiff(arVariantsP[i]) < GetMmrDiff(arVariantsP[i - 1])) {
					arVariantsP[i], arVariantsP[i - 1] = arVariantsP[i - 1], arVariantsP[i]; //switch
					if (bSorted) {
						bSorted = false;
					}
				}
				if (!bSorted) {
					for i := iSize - 2; i >= 0; i-- {
						if (GetMmrDiff(arVariantsP[i]) > GetMmrDiff(arVariantsP[i + 1])) {
							arVariantsP[i], arVariantsP[i + 1] = arVariantsP[i + 1], arVariantsP[i]; //switch
						}
					}
				}
			}
		}
	}

	//Cut them by mmr diff of 500
	var arTryCut [][2][]*players.EntPlayer;
	for i, _ := range arVariantsP {
		if (GetMmrDiff(arVariantsP[i]) >= int(settings.MmrDiffGuaranteedWin) / 2) {
			arTryCut = arVariantsP[:i];
			break;
		}
	}
	if (len(arTryCut) > 0) {
		arVariantsP = arTryCut;
	} else {
		//If none left, get the first variant from previous step, and return it.
		return arVariantsP[0][0], arVariantsP[0][1];
	}

	//find variants with the biggest amount of happy duos, select randomly
	var arMostOfHappyDuos []int;
	var iCurMostOfHappyDuos int;
	for i, _ := range arVariantsP {
		iHappyDuos := GetHappyDuoPlayers(arVariantsP[i]);
		if (iHappyDuos > iCurMostOfHappyDuos) {
			iCurMostOfHappyDuos = iHappyDuos;
			arMostOfHappyDuos = make([]int, 0);
		}
		if (iHappyDuos == iCurMostOfHappyDuos) {
			arMostOfHappyDuos = append(arMostOfHappyDuos, i);
		}
	}

	iRandInt, _ := utils.GetRandInt(0, len(arMostOfHappyDuos) - 1);
	return arVariantsP[arMostOfHappyDuos[iRandInt]][0], arVariantsP[arMostOfHappyDuos[iRandInt]][1];
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

func GetMmrDiff(arVariantP [2][]*players.EntPlayer) int {
	var iMmr [2]int;
	for iT := 0; iT < 2; iT++ {
		for iP := 0; iP < 4; iP++ {
			iMmr[iT] = iMmr[iT] + arVariantP[iT][iP].Mmr;
		}
	}
	iDiff := (iMmr[0] - iMmr[1]) / 4;
	if (iDiff < 0) {
		return iDiff * -1;
	}
	return iDiff;
}