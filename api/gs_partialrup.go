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

			games.MuGames.Lock();
			pGame := games.GetGameByIP(sIP);
			if (pGame != nil) {

				sReadyPlayers := c.PostForm("ready_players");
				if (sReadyPlayers != "") {

					sResponse = fmt.Sprintf("%s\n	\"success\" \"1\"", sResponse);

					select {
					case pGame.ReceiverReadyList <- strings.Split(sReadyPlayers, ","):
					default:
					}

				} else {
					sResponse = fmt.Sprintf("%s\n	\"success\" \"0\"", sResponse);
					sResponse = fmt.Sprintf("%s\n	\"error\" \"Bad ready_players parameter\"", sResponse);
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
