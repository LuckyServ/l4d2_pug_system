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


func HttpReqGSChatLogs(c *gin.Context) {

	var sResponse string = "\"VDFresponse\"\n{";

	sAuthKey := c.PostForm("auth_key");
	sLogLine := c.PostForm("logline");
	sSteamID64 := c.PostForm("steamid64");
	if (auth.Backend(sAuthKey)) {
		sIP := c.PostForm("ip");
		if (sIP != "" && sLogLine != "" && sSteamID64 != "") {
			games.MuGames.RLock();
			pGame := games.GetGameByIP(sIP);
			if (pGame != nil) {
				sGameID := pGame.ID;
				games.MuGames.RUnlock();

				sResponse = fmt.Sprintf("%s\n	\"success\" \"1\"", sResponse);

				go database.GameServerChatLog(time.Now().UnixMilli(), sGameID, sSteamID64, base64.StdEncoding.EncodeToString([]byte(sLogLine)));

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
