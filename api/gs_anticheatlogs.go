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


func HttpReqGSAntiCheatLogs(c *gin.Context) {

	var sResponse string = "\"VDFresponse\"\n{";

	sAuthKey := c.PostForm("auth_key");
	sLogLine := c.PostForm("logline");
	if (auth.Backend(sAuthKey)) {
		sIP := c.PostForm("ip");
		if (sIP != "" && sLogLine != "") {
			var sGameID string;
			games.MuGames.RLock();
			pGame := games.GetGameByIP(sIP);
			if (pGame != nil) {
				sGameID = pGame.ID;
			} else {
				sGameID = "Not a l4d2center game";
			}
			games.MuGames.RUnlock();


			sResponse = fmt.Sprintf("%s\n	\"success\" \"1\"", sResponse);

			sBuffer := fmt.Sprintf("Time: %s, game %s, log line: %s", time.Now().Format("02/Jan/2006 - 15:04:05.00 MST"), sGameID, sLogLine);
			go database.AntiCheatLog(base64.StdEncoding.EncodeToString([]byte(sBuffer)));

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
