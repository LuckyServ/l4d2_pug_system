package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../games"
)


func HttpReqGSFullReadyUp(c *gin.Context) {

	var sResponse string = "\"VDFresponse\"\n{";


	sAuthKey := c.PostForm("auth_key");
	if (auth.Backend(sAuthKey)) {
		sIP := c.PostForm("ip");
		if (sIP != "") {

			games.MuGames.RLock();
			pGame := games.GetGameByIP(sIP);
			if (pGame != nil) {

				sResponse = fmt.Sprintf("%s\n	\"success\" \"1\"", sResponse);

				select {
				case pGame.ReceiverFullRUP <- true:
				default:
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
