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
	SomeoneBanned		bool
}


func UpdateMmr(oResult EntGameResult, arPlayers [2][]*players.EntPlayer) { //Players must be locked outside

	if (oResult.SettledScores[0] == 0 && oResult.SettledScores[1] == 0) { //dont touch mmr if scores are zero
		return;
	} else if (oResult.SomeoneBanned) { //dont touch mmr if game stopped due to a ban
		return;
	}

	//Get winner and winning coef
	var f32WinCoef float32;
	var iWinner int = -1; //-1 == draw
	if (oResult.SettledScores[0] > oResult.SettledScores[1]) {
		iWinner = 0;
		f32WinCoef = (float32(oResult.SettledScores[0]) / (float32(oResult.SettledScores[1]) + 0.001/*in case if scores are 0*/)) - 1.0;
	} else if (oResult.SettledScores[1] > oResult.SettledScores[0]) {
		iWinner = 1;
		f32WinCoef = (float32(oResult.SettledScores[1]) / (float32(oResult.SettledScores[0]) + 0.001/*in case if scores are 0*/)) - 1.0;
	}
	//twice as high scores of a winner = absolute win
	if (f32WinCoef > 1.0) {
		f32WinCoef = 1.0;
	}
	//Prevent too high mmr gains (and losses) if played 1 map only
	if (oResult.MapsFinished == 1 && f32WinCoef > 0.0) {
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
		f32FavorCoef = ((f32TeamMmr[0] - f32TeamMmr[1]) / 4) / settings.MmrDiffGuaranteedWin;
	} else if (f32TeamMmr[1] > f32TeamMmr[0]) {
		iFavorited = 1;
		f32FavorCoef = ((f32TeamMmr[1] - f32TeamMmr[0]) / 4) / settings.MmrDiffGuaranteedWin;
	}
	if (f32FavorCoef > 1.0) {
		f32FavorCoef = 1.0;
	}

	//How much mmr team A gets (can be negative)
	var f32TeamAgets float32;
	if (iWinner == iFavorited) {
		if (iWinner != -1) {
			f32TeamAgets = ((settings.MmrAbsoluteWin - settings.MmrMinimumWin) * (f32WinCoef - f32FavorCoef)) + settings.MmrMinimumWin;
		}
	} else {
		if (iWinner != -1) {
			f32TeamAgets = ((settings.MmrAbsoluteWin - settings.MmrMinimumWin) * (f32WinCoef + f32FavorCoef)) + settings.MmrMinimumWin;
		} else {
			f32TeamAgets = (settings.MmrAbsoluteWin - settings.MmrMinimumWin) * f32FavorCoef;
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

	//check what team did RQ
	var iRQTeam int = -1; //A (0) or B (1) or both (2), or nonr (-1)
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

	//Apply new mmr
	iTeamAgets := int(math.Round(float64(f32TeamAgets)));
	if (iTeamAgets > 0) {
		for iP := 0; iP < 4; iP++ {
			arPlayers[0][iP].LastGameResult = 3; //won
			if (utils.GetStringIdxInArray(arPlayers[0][iP].SteamID64, oResult.AbsentPlayers) != -1) {
				continue;
			} else if (iRQTeam == 0 && utils.GetDifferenceInt(oResult.SettledScores[0], oResult.SettledScores[1]) < 800) {
				continue;
			}
			arPlayers[0][iP].Mmr = arPlayers[0][iP].Mmr + (iTeamAgets + int(f32TeamAgets * arPlayers[0][iP].MmrUncertainty));
			arPlayers[0][iP].MmrUncertainty = arPlayers[0][iP].MmrUncertainty * 0.8; //reduce uncertainty
		}
		for iP := 0; iP < 4; iP++ {
			arPlayers[1][iP].LastGameResult = 2; //lost
			if (iRQTeam == 0 && utils.GetDifferenceInt(oResult.SettledScores[0], oResult.SettledScores[1]) < 800) {
				continue;
			}
			arPlayers[1][iP].Mmr = arPlayers[1][iP].Mmr - (iTeamAgets + int(f32TeamAgets * arPlayers[1][iP].MmrUncertainty));
			arPlayers[1][iP].MmrUncertainty = arPlayers[1][iP].MmrUncertainty * 0.8; //reduce uncertainty
		}
	} else if (iTeamAgets < 0) {
		for iP := 0; iP < 4; iP++ {
			arPlayers[0][iP].LastGameResult = 2; //lost
			if (iRQTeam == 1 && utils.GetDifferenceInt(oResult.SettledScores[0], oResult.SettledScores[1]) < 800) {
				continue;
			}
			arPlayers[0][iP].Mmr = arPlayers[0][iP].Mmr + (iTeamAgets + int(f32TeamAgets * arPlayers[0][iP].MmrUncertainty));
			arPlayers[0][iP].MmrUncertainty = arPlayers[0][iP].MmrUncertainty * 0.8; //reduce uncertainty
		}
		for iP := 0; iP < 4; iP++ {
			arPlayers[1][iP].LastGameResult = 3; //won
			if (utils.GetStringIdxInArray(arPlayers[1][iP].SteamID64, oResult.AbsentPlayers) != -1) {
				continue;
			} else if (iRQTeam == 1 && utils.GetDifferenceInt(oResult.SettledScores[0], oResult.SettledScores[1]) < 800) {
				continue;
			}
			arPlayers[1][iP].Mmr = arPlayers[1][iP].Mmr - (iTeamAgets + int(f32TeamAgets * arPlayers[1][iP].MmrUncertainty));
			arPlayers[1][iP].MmrUncertainty = arPlayers[1][iP].MmrUncertainty * 0.8; //reduce uncertainty
		}
	} else { //draw
		for iT := 0; iT < 2; iT++ {
			for iP := 0; iP < 4; iP++ {
				arPlayers[iT][iP].MmrUncertainty = arPlayers[iT][iP].MmrUncertainty * 0.8; //reduce uncertainty
				arPlayers[iT][iP].LastGameResult = 1; //draw
			}
		}
	}
}


/*func DetermineFinalScores(oResult EntGameResult, arPlayers [2][]*players.EntPlayer) [2]int { //Players must be locked outside
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

		//everything else
		arScoresBuffer := oResult.SettledScores;
		arScoresBuffer[iRQTeam] = utils.MaxValInt(arScoresBuffer[iRQTeam] - settings.RQDefaultPenalty, 0);
		return arScoresBuffer;

	} else {
		//case: game ended with no RQs
		return oResult.SettledScores;
	}
	return [2]int{0, 0};
}*/
