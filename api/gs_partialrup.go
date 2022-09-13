package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"../games"
	"../players/auth"
)


func HttpReqGSPartialReadyUp(c *gin.Context) {

	var sResponse string = "\"VDFresponse\"\n{";


	sAuthKey := c.PostForm("auth_key");
	if (auth.Backend(sAuthKey)) {
		sIP := c.PostForm("ip");
		if (sIP != "") {

			games.MuGames.RLock();
			pGame := games.GetGameByIP(sIP);
			if (pGame != nil) {

				sReadyPlayers := c.PostForm("ready_players");
				if (sReadyPlayers != "") {

					select {
					case pGame.ReceiverReadyList <- strings.Split(sReadyPlayers, ","):
						sResponse = fmt.Sprintf("%s\n	\"success\" \"1\"", sResponse);
					default:
						sResponse = fmt.Sprintf("%s\n	\"success\" \"0\"", sResponse);
						sResponse = fmt.Sprintf("%s\n	\"error\" \"Not waiting for ready players list\"", sResponse);
					}

				} else {
					sResponse = fmt.Sprintf("%s\n	\"success\" \"0\"", sResponse);
					sResponse = fmt.Sprintf("%s\n	\"error\" \"Bad ready_players parameter\"", sResponse);
				}

			} else {
				sResponse = fmt.Sprintf("%s\n	\"success\" \"0\"", sResponse);
				sResponse = fmt.Sprintf("%s\n	\"error\" \"No game on this IP\"", sResponse);
			}
			games.MuGames.RUnlock();
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
