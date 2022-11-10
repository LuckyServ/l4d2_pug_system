package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../games"
	"../players"
	"../settings"
)


func HttpReqGSGetGame(c *gin.Context) {

	var sResponse string = "\"VDFresponse\"\n{";


	sAuthKey := c.PostForm("auth_key");
	if (auth.Backend(sAuthKey)) {
		sIP := c.PostForm("ip");
		if (sIP != "") {
			players.MuPlayers.RLock();
			games.MuGames.RLock();
			pGame := games.GetGameByIP(sIP);
			if (pGame != nil) {

				sResponse = fmt.Sprintf("%s\n	\"success\" \"1\"", sResponse);
				sResponse = fmt.Sprintf("%s\n	\"game_id\" \"%s\"", sResponse, pGame.ID);

				for i := 0; i < 4; i++ {
					sResponse = fmt.Sprintf("%s\n	\"player_a%d\" \"%s\"", sResponse, i, pGame.PlayersA[i].SteamID64 + "s"); //"s" is a workaround for sourcemod bug, it recognizes the value as int otherwise
				}
				for i := 0; i < 4; i++ {
					sResponse = fmt.Sprintf("%s\n	\"player_b%d\" \"%s\"", sResponse, i, pGame.PlayersB[i].SteamID64 + "s"); //"s" is a workaround for sourcemod bug, it recognizes the value as int otherwise
				}

				sResponse = fmt.Sprintf("%s\n	\"confogl\" \"%s\"", sResponse, pGame.GameConfig.CodeName);
				sResponse = fmt.Sprintf("%s\n	\"first_map\" \"%s\"", sResponse, pGame.Maps[0]);
				sResponse = fmt.Sprintf("%s\n	\"last_map\" \"%s\"", sResponse, pGame.Maps[len(pGame.Maps) - 1]);
				sResponse = fmt.Sprintf("%s\n	\"max_absent\" \"%d\"", sResponse, settings.MaxAbsentSeconds);
				sResponse = fmt.Sprintf("%s\n	\"max_single_absent\" \"%d\"", sResponse, settings.MaxSingleAbsentSeconds);

				if (pGame.State == games.StateWaitPlayersJoin) {
					sResponse = fmt.Sprintf("%s\n	\"game_state\" \"wait_readyup\"", sResponse);
				} else if (pGame.State == games.StateReadyUpExpired) {
					sResponse = fmt.Sprintf("%s\n	\"game_state\" \"readyup_expired\"", sResponse);
				} else if (pGame.State == games.StateGameProceeds) {
					sResponse = fmt.Sprintf("%s\n	\"game_state\" \"game_proceeds\"", sResponse);
				} else {
					sResponse = fmt.Sprintf("%s\n	\"game_state\" \"other\"", sResponse);
				}

			} else {
				sResponse = fmt.Sprintf("%s\n	\"success\" \"0\"", sResponse);
				sResponse = fmt.Sprintf("%s\n	\"error\" \"No game on this IP\"", sResponse);
			}
			games.MuGames.RUnlock();
			players.MuPlayers.RUnlock();
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
