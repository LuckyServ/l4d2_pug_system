package rating

import (
	"../players"
	"../utils"
	"../settings"
	"math"
)

type EntGameResult struct {
	SettledScores		[2]int
	CurrentScores		[2]int
	InRound				bool
	CurrentHalf			int //1 or 2
	TeamsFlipped		bool
	TankKilled			bool //only valid if InRound == true
	TankInPlay			bool //only valid if InRound == true
	Dominator			[2]string
	Inferior			[2]string
	GameEnded			bool //no more results should be accepted
	InMapTransition		bool
	IsLastMap			bool
	AbsentPlayers		[]string
	ConnectedPlayers	int
	MapsFinished		int
}


func UpdateMmr(oResult EntGameResult, arFinalScores [2]int, arPlayers [2][]*players.EntPlayer) { //Players must be locked outside

	//Get winner and winning coef
	var f32WinCoef float32;
	var iWinner int = -1; //-1 == draw
	if (arFinalScores[0] > arFinalScores[1]) {
		iWinner = 0;
		f32WinCoef = (float32(arFinalScores[0]) / (float32(arFinalScores[1]) + 0.001/*in case if scores are 0*/)) - 1.0;
	} else if (arFinalScores[1] > arFinalScores[0]) {
		iWinner = 1;
		f32WinCoef = (float32(arFinalScores[1]) / (float32(arFinalScores[0]) + 0.001/*in case if scores are 0*/)) - 1.0;
	}
	//twice as high scores of a winner = absolute win
	if (f32WinCoef > 1.0) {
		f32WinCoef = 1.0;
	}
	//Prevent too high mmr gains (and losses) if played 1 map only
	if ((oResult.MapsFinished == 0 || (oResult.MapsFinished == 1 && (oResult.InMapTransition || (oResult.CurrentHalf == 1 && oResult.InRound)))) && f32WinCoef > 0.0) {
		f32WinCoef = 0.0;
	}

	//Get favorited to win
	var f32FavorCoef float32;
	var iFavorited int = -1; //-1 == equal mmr
	var f32TeamMmr [2]float32;
	for iT := 0; iT < 2; iT++ {
		f32TeamMmr[iT] = float32(arPlayers[iT][0].Mmr + arPlayers[iT][1].Mmr + arPlayers[iT][2].Mmr + arPlayers[iT][3].Mmr);
	}
	if (f32TeamMmr[0] > f32TeamMmr[1]) {
		iFavorited = 0;
		f32FavorCoef = (f32TeamMmr[0] / (f32TeamMmr[1] + 0.001/*in case if all mmr are 0*/)) - 1.0;
	} else if (f32TeamMmr[1] > f32TeamMmr[0]) {
		iFavorited = 1;
		f32FavorCoef = (f32TeamMmr[1] / (f32TeamMmr[0] + 0.001/*in case if all mmr are 0*/)) - 1.0;
	}
	if (f32FavorCoef > 1.0) {
		f32FavorCoef = 1.0;
	}

	//How much mmr team A gets (can be negative)
	var f32TeamAgets float32;
	if (iWinner == iFavorited) {
		if (iWinner != -1) {
			f32TeamAgets = (settings.MmrAbsoluteWin * (f32WinCoef - f32FavorCoef)) + settings.MmrMinimumWin;
		}
	} else {
		if (iWinner != -1) {
			f32TeamAgets = (settings.MmrAbsoluteWin * (f32WinCoef + f32FavorCoef)) + settings.MmrMinimumWin;
		} else {
			f32TeamAgets = settings.MmrAbsoluteWin * f32FavorCoef;
			if (iFavorited == 0) {
				f32TeamAgets = f32TeamAgets * -1.0;
			}
		}
	}
	if (iWinner == 1) {
		f32TeamAgets = f32TeamAgets * -1.0;
	}
	//limit max mmr gain
	if (f32TeamAgets > settings.MmrAbsoluteWin) {
		f32TeamAgets = settings.MmrAbsoluteWin;
	} else if (f32TeamAgets < settings.MmrAbsoluteWin * -1.0) {
		f32TeamAgets = settings.MmrAbsoluteWin * -1.0;
	}

	//Apply new team based mmr
	iTeamAgets := int(math.Round(float64(f32TeamAgets)));
	if (iTeamAgets > 0) {
		for iP := 0; iP < 4; iP++ {
			if (arPlayers[0][iP].SteamID64 != oResult.Inferior[0]) {
				arPlayers[0][iP].Mmr = arPlayers[0][iP].Mmr + iTeamAgets;
			}
		}
		for iP := 0; iP < 4; iP++ {
			if (arPlayers[1][iP].SteamID64 != oResult.Dominator[1]) {
				arPlayers[1][iP].Mmr = arPlayers[1][iP].Mmr - iTeamAgets;
			}
		}
	} else if (iTeamAgets < 0) {
		for iP := 0; iP < 4; iP++ {
			if (arPlayers[0][iP].SteamID64 != oResult.Dominator[0]) {
				arPlayers[0][iP].Mmr = arPlayers[0][iP].Mmr + iTeamAgets;
			}
		}
		for iP := 0; iP < 4; iP++ {
			if (arPlayers[1][iP].SteamID64 != oResult.Inferior[1]) {
				arPlayers[1][iP].Mmr = arPlayers[1][iP].Mmr - iTeamAgets;
			}
		}
	}

	//Shift mmr to the positive side
	var iMmrShift int; //we dont allow mmr lower than 1
	for iT := 0; iT < 2; iT++ {
		for iP := 0; iP < 4; iP++ {
			if (arPlayers[iT][iP].Mmr < 1 && iMmrShift < 1 - arPlayers[iT][iP].Mmr) {
				iMmrShift = 1 - arPlayers[iT][iP].Mmr;
			}
		}
	}
	for iT := 0; iT < 2; iT++ {
		for iP := 0; iP < 4; iP++ {
			arPlayers[iT][iP].Mmr = arPlayers[iT][iP].Mmr + iMmrShift;
		}
	}

	//Also take into account mmr deviation (uncertainty)
	if (iTeamAgets > 0.0) {
		for iP := 0; iP < 4; iP++ {
			if (arPlayers[0][iP].SteamID64 != oResult.Inferior[0]) {
				arPlayers[0][iP].Mmr = arPlayers[0][iP].Mmr + int(f32TeamAgets * arPlayers[0][iP].MmrUncertainty);
			}
		}
		for iP := 0; iP < 4; iP++ {
			if (arPlayers[1][iP].SteamID64 != oResult.Dominator[1]) {
				arPlayers[1][iP].Mmr = arPlayers[1][iP].Mmr - int(f32TeamAgets * arPlayers[1][iP].MmrUncertainty);
			}
		}
	} else if (iTeamAgets < 0.0) {
		for iP := 0; iP < 4; iP++ {
			if (arPlayers[0][iP].SteamID64 != oResult.Dominator[0]) {
				arPlayers[0][iP].Mmr = arPlayers[0][iP].Mmr + int(f32TeamAgets * arPlayers[0][iP].MmrUncertainty);
			}
		}
		for iP := 0; iP < 4; iP++ {
			if (arPlayers[1][iP].SteamID64 != oResult.Inferior[1]) {
				arPlayers[1][iP].Mmr = arPlayers[1][iP].Mmr - int(f32TeamAgets * arPlayers[1][iP].MmrUncertainty);
			}
		}
	}

	//Cut all mmr's below 1
	for iT := 0; iT < 2; iT++ {
		for iP := 0; iP < 4; iP++ {
			if (arPlayers[iT][iP].Mmr < 1) {
				arPlayers[iT][iP].Mmr = 1;
			}
			arPlayers[iT][iP].MmrUncertainty = arPlayers[iT][iP].MmrUncertainty * 0.75; //reduce uncertainty
		}
	}
}


func DetermineFinalScores(oResult EntGameResult, arPlayers [2][]*players.EntPlayer) [2]int { //Players must be locked outside
	if (len(oResult.AbsentPlayers) > 0) {

		//check what team did RQ (whole team gets responsible for one player RQ)
		var iRQTeam int = -1; //A (0) or B (1) or both (2)
		for _, pPlayer := range arPlayers[0] {
			if (utils.GetStringIdxInArray(pPlayer.SteamID64, oResult.AbsentPlayers) != -1) {
				iRQTeam = 0;
				break;
			}
		}
		for _, pPlayer := range arPlayers[1] {
			if (utils.GetStringIdxInArray(pPlayer.SteamID64, oResult.AbsentPlayers) != -1) {
				if (iRQTeam == -1) {
					iRQTeam = 1;
					break;
				} else {
					iRQTeam = 2;
					break;
				}
			}
		}

		//case: RQd players are > 1 and from different teams
		if (iRQTeam != 0 && iRQTeam != 1) {
			return oResult.SettledScores;
		}

		//case: a player RQd om map change
		if (oResult.InMapTransition) {
			arScoresBuffer := oResult.SettledScores;
			arScoresBuffer[iRQTeam] = utils.MaxValInt(arScoresBuffer[iRQTeam] - settings.RQMidMapTransPenalty, 0);
			return arScoresBuffer;
		}

		//case: a losing player RQd on side switch
		arChaptScores := [2]int{oResult.CurrentScores[0] - oResult.SettledScores[0], oResult.CurrentScores[1] - oResult.SettledScores[1]};
		if (!oResult.InRound && !oResult.InMapTransition) {
			var iLosingTeam int = -1;
			if (arChaptScores[0] > 0 && arChaptScores[1] == 0) {
				iLosingTeam = 1;
			} else if (arChaptScores[0] == 0 && arChaptScores[1] > 0) {
				iLosingTeam = 0;
			} else { //weird situation
				return oResult.SettledScores;
			}

			if (iLosingTeam == iRQTeam) {
				return oResult.CurrentScores;
			} else {
				return oResult.SettledScores;
			}
		}

		//case: infected player left midtank or after killing tank, on 2nd half
		if (oResult.InRound && oResult.CurrentHalf == 2 && ((iRQTeam == 0 && oResult.TeamsFlipped) || (iRQTeam == 1 && !oResult.TeamsFlipped)) && (oResult.TankInPlay || oResult.TankKilled)) {
			arScoresBuffer := oResult.CurrentScores;
			arScoresBuffer[iRQTeam] = utils.MaxValInt(arScoresBuffer[iRQTeam] - settings.RQInfHalf2MidTank, 0);
			return arScoresBuffer;
		}

		//case: survivor player left midgame on 2nd half
		if (oResult.InRound && oResult.CurrentHalf == 2 && ((iRQTeam == 0 && !oResult.TeamsFlipped) || (iRQTeam == 1 && oResult.TeamsFlipped))) {
			return oResult.CurrentScores;
		}

		//Anything else
		return oResult.SettledScores;

	} else {
		//case: game ended with no RQs
		return oResult.SettledScores;
	}
	return [2]int{0, 0};
}
