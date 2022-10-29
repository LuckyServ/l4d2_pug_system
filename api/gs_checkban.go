package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"../players"
	"../bans"
)


func HttpReqGSCheckBan(c *gin.Context) {

	var sResponse string = "\"VDFresponse\"\n{";

	sSteamID64 := c.PostForm("steamid64");

	if (sSteamID64 != "") {
		sResponse = fmt.Sprintf("%s\n	\"success\" \"1\"", sResponse);
		sResponse = fmt.Sprintf("%s\n	\"steamid64\" \"%s\"", sResponse, sSteamID64 + "s");

		players.MuPlayers.RLock();
		pPlayer, bFound := players.MapPlayers[sSteamID64];
		if (bFound) {
			if (pPlayer.Access <= -2) {
				sResponse = fmt.Sprintf("%s\n	\"isbanned\" \"1\"", sResponse);
			} else {
				sResponse = fmt.Sprintf("%s\n	\"isbanned\" \"0\"", sResponse);
			}
			players.MuPlayers.RUnlock();
		} else {
			players.MuPlayers.RUnlock();


			var bIsDatabaseBanned bool;
			bans.ChanLock <- true;
			iBanlistSize := len(bans.ArrayBanRecords);
			for i := iBanlistSize - 1; i >= 0; i-- {
				if (bans.ArrayBanRecords[i].AcceptedAt == 0 && bans.ArrayBanRecords[i].SteamID64 == sSteamID64) {
					bIsDatabaseBanned = true;
					break;
				}
			}
			bans.ChanUnlock <- true;

			if (bIsDatabaseBanned) {
				sResponse = fmt.Sprintf("%s\n	\"isbanned\" \"1\"", sResponse);
			} else {
				sResponse = fmt.Sprintf("%s\n	\"isbanned\" \"0\"", sResponse);
			}

		}

	} else {
		sResponse = fmt.Sprintf("%s\n	\"success\" \"0\"", sResponse);
		sResponse = fmt.Sprintf("%s\n	\"error\" \"Bad parameters\"", sResponse);
	}


	sResponse = sResponse + "\n}\n";
	c.String(200, sResponse);
}
