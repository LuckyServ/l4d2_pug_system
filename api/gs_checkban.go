package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"../players"
	"../bans"
	"strings"
	"encoding/base64"
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
				byBanReason, _ := base64.StdEncoding.DecodeString(pPlayer.BanReason);
				sResponse = fmt.Sprintf("%s\n	\"ban_reason\" \"%s\"", sResponse, strings.ReplaceAll(string(byBanReason), `"`, `\"`));
			} else {
				sResponse = fmt.Sprintf("%s\n	\"isbanned\" \"0\"", sResponse);
			}
			players.MuPlayers.RUnlock();
		} else {
			players.MuPlayers.RUnlock();


			var sDatabaseBanReason string;
			bans.ChanLock <- true;
			iBanlistSize := len(bans.ArrayBanRecords);
			for i := iBanlistSize - 1; i >= 0; i-- {
				if (bans.ArrayBanRecords[i].AcceptedAt == 0 && bans.ArrayBanRecords[i].SteamID64 == sSteamID64) {
					byBanReason, _ := base64.StdEncoding.DecodeString(bans.ArrayBanRecords[i].BanReasonBase64);
					sDatabaseBanReason = string(byBanReason);
					break;
				}
			}
			bans.ChanUnlock <- true;

			if (sDatabaseBanReason != "") {
				sResponse = fmt.Sprintf("%s\n	\"isbanned\" \"1\"", sResponse);
				sResponse = fmt.Sprintf("%s\n	\"ban_reason\" \"%s\"", sResponse, strings.ReplaceAll(sDatabaseBanReason, `"`, `\"`));
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
