package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"../games"
	"../settings"
	"../players/auth"
	"../players"
)


func HttpReqGSGameResults(c *gin.Context) {

	var sResponse string = "\"VDFresponse\"\n{";

	sAuthKey := c.PostForm("auth_key");
	if (auth.Backend(sAuthKey)) {
		sIP := c.PostForm("ip");
		if (sIP != "") {
			games.MuGames.Lock();
			pGame := games.GetGameByIP(sIP);
			if (pGame != nil) {

				//we wont check for errors
				iSettledScoresA, _ := strconv.Atoi(c.PostForm("settled_scores_a"));
				iSettledScoresB, _ := strconv.Atoi(c.PostForm("settled_scores_b"));
				iCurrentScoresA, _ := strconv.Atoi(c.PostForm("current_scores_a"));
				iCurrentScoresB, _ := strconv.Atoi(c.PostForm("current_scores_b"));
				bInRound := (c.PostForm("in_round") == "yes");
				iHalf, _ := strconv.Atoi(c.PostForm("half"));
				bTeamsFlipped := (c.PostForm("teams_flipped") == "yes");
				bTankKilled := (c.PostForm("tank_killed") == "yes");
				sDominatorA := c.PostForm("dominator_a");
				sDominatorB := c.PostForm("dominator_b");
				sInferiorA := c.PostForm("inferior_a");
				sInferiorB := c.PostForm("inferior_b");
				bGameFinished := (c.PostForm("game_finished") == "yes");
				iPlayers, _ := strconv.Atoi(c.PostForm("players_connected"));

				players.MuPlayers.Lock();

				var arAbsentPlayers []string;
				for _, pPlayer := range pGame.PlayersUnpaired {
					iAbsentFor, errAbsentFor := strconv.Atoi(c.PostForm(pPlayer.SteamID64));
					if (errAbsentFor == nil && int64(iAbsentFor) >= settings.MaxAbsentSeconds) {
						arAbsentPlayers = append(arAbsentPlayers, pPlayer.SteamID64);
					}
				}

				players.MuPlayers.Unlock();


				oResult := games.EntGameResult{
					SettledScores:			[2]int{iSettledScoresA, iSettledScoresB},
					CurrentScores:			[2]int{iCurrentScoresA, iCurrentScoresB},
					InRound:				bInRound,
					CurrentHalf:			iHalf,
					TeamsFlipped:			bTeamsFlipped,
					TankKilled:				bTankKilled,
					Dominator:				[2]string{sDominatorA, sDominatorB},
					Inferior:				[2]string{sInferiorA, sInferiorB},
					GameEnded:				(bGameFinished || len(arAbsentPlayers) > 0 || iPlayers <= settings.MinPlayersCount),
					AbsentPlayers:			arAbsentPlayers,
					ConnectedPlayers:		iPlayers,
				};

				select {
				case pGame.ReceiverResult <- oResult:
					sResponse = fmt.Sprintf("%s\n	\"success\" \"1\"", sResponse);
					if (bGameFinished) {
						sResponse = fmt.Sprintf("%s\n	\"game_ended_type\" \"1\"", sResponse);
					} else if (len(arAbsentPlayers) > 0) {
						sResponse = fmt.Sprintf("%s\n	\"game_ended_type\" \"2\"", sResponse);
					} else if (iPlayers <= settings.MinPlayersCount) {
						sResponse = fmt.Sprintf("%s\n	\"game_ended_type\" \"3\"", sResponse);
					} else {
						sResponse = fmt.Sprintf("%s\n	\"game_ended_type\" \"0\"", sResponse);
					}
				default:
					sResponse = fmt.Sprintf("%s\n	\"success\" \"0\"", sResponse);
					sResponse = fmt.Sprintf("%s\n	\"error\" \"Not waiting for game results\"", sResponse);
				}



			} else {
				sResponse = fmt.Sprintf("%s\n	\"success\" \"0\"", sResponse);
				sResponse = fmt.Sprintf("%s\n	\"error\" \"No game on this IP\"", sResponse);
			}
			games.MuGames.Unlock();
		} else {
			sResponse = fmt.Sprintf("%s\n	\"success\" \"0\"", sResponse);
			sResponse = fmt.Sprintf("%s\n	\"error\" \"No ip parameter\"", sResponse);
		}
	} else {
		sResponse = fmt.Sprintf("%s\n	\"success\" \"0\"", sResponse);
		sResponse = fmt.Sprintf("%s\n	\"error\" \"Bad auth key\"", sResponse);
	}

	sResponse = sResponse + "\n}\n";
	c.String(200, sResponse);
}
