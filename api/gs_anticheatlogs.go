package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../games"
	"../database"
	"encoding/base64"
	"time"
)


func HttpReqAntiCheatLogs(c *gin.Context) {

	var sResponse string = "\"VDFresponse\"\n{";

	sAuthKey := c.PostForm("auth_key");
	sLogLine := c.PostForm("logline");
	if (auth.Backend(sAuthKey)) {
		sIP := c.PostForm("ip");
		if (sIP != "" && sLogLine != "") {
			games.MuGames.RLock();
			pGame := games.GetGameByIP(sIP);
			if (pGame != nil) {
				sGameID := pGame.ID;
				games.MuGames.RUnlock();

				sResponse = fmt.Sprintf("%s\n	\"success\" \"1\"", sResponse);

				sBuffer := fmt.Sprintf("Time: %s, game %s, log line: %s", time.Now().Format("01/Jan/2006 - 15:04:05.00 MST"), sGameID, sLogLine);
				go database.AntiCheatLog(base64.StdEncoding.EncodeToString([]byte(sBuffer)));

			} else {
				games.MuGames.RUnlock();
				sResponse = fmt.Sprintf("%s\n	\"success\" \"0\"", sResponse);
				sResponse = fmt.Sprintf("%s\n	\"error\" \"No game on this IP\"", sResponse);
			}
		} else {
			sResponse = fmt.Sprintf("%s\n	\"success\" \"0\"", sResponse);
			sResponse = fmt.Sprintf("%s\n	\"error\" \"Bad parameters\"", sResponse);
		}
	} else {
		sResponse = fmt.Sprintf("%s\n	\"success\" \"0\"", sResponse);
		sResponse = fmt.Sprintf("%s\n	\"error\" \"Bad auth key\"", sResponse);
	}

	sResponse = sResponse + "\n}\n";
	c.String(200, sResponse);
}
