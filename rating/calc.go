package rating

import (
	"../players"
	"../utils"
	"../settings"
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


func DetermineFinalScores(oResult EntGameResult, arPlayersA []*players.EntPlayer, arPlayersB []*players.EntPlayer) [2]int { //Players must be locked outside
	if (len(oResult.AbsentPlayers) > 0) {

		//check what team did RQ (whole team gets responsible for one player RQ)
		var iRQTeam int = -1; //A (0) or B (1) or both (2)
		for _, pPlayer := range arPlayersA {
			if (utils.GetStringIdxInArray(pPlayer.SteamID64, oResult.AbsentPlayers) != -1) {
				iRQTeam = 0;
				break;
			}
		}
		for _, pPlayer := range arPlayersB {
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