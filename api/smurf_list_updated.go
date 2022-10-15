package api

import (
	"github.com/gin-gonic/gin"
	"strings"
	"../bans"
	"../players"
	"../players/auth"
	"fmt"
	"time"
	"encoding/base64"
)



func HttpReqSMURFListUpdated(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	mapResponse["success"] = false;
	sAuthKey := c.Query("auth_key");
	if (auth.Backend(sAuthKey)) {
		sAccounts := c.Query("accounts");
		if (sAccounts != "") {
			mapResponse["success"] = true;
			arAccounts := strings.Split(sAccounts, ",");
			if (len(arAccounts) > 1) {
				//check if any one of them is banned
				bans.ChanLock <- true;

				var oMimicBanRecord bans.EntBanRecord;
				i64CurTime := time.Now().UnixMilli();

				iBanlistSize := len(bans.ArrayBanRecords);
				for _, sSteamID64 := range arAccounts {
					for i := iBanlistSize - 1; i >= 0; i-- {
						if (bans.ArrayBanRecords[i].SteamID64 == sSteamID64) {
							if (bans.ArrayBanRecords[i].AcceptedAt == 0 || bans.ArrayBanRecords[i].AcceptedAt + bans.ArrayBanRecords[i].BanLength > i64CurTime) {
								oMimicBanRecord = bans.ArrayBanRecords[i];
								break;
							}
						}
					}
					if (oMimicBanRecord.CreatedAt > 0) {
						break;
					}
				}

				bans.ChanUnlock <- true;

				if (oMimicBanRecord.CreatedAt > 0) {
					//check if any one of them is online || in game || in lobby
					var arToBan []string;
					players.MuPlayers.RLock();
					for _, sSteamID64 := range arAccounts {
						pPlayer, bFound := players.MapPlayers[sSteamID64];
						if (bFound && (pPlayer.IsOnline || pPlayer.IsInGame || pPlayer.IsInLobby) && pPlayer.Access >= -1) {
							arToBan = append(arToBan, sSteamID64);
						}
					}
					players.MuPlayers.RUnlock();


					//ban those found on previous step
					for _, sSteamID64 := range arToBan {
						byNickname, _ := base64.StdEncoding.DecodeString(oMimicBanRecord.NicknameBase64);
						oBanReq := bans.EntManualBanReq{
							SteamID64:			sSteamID64,
							Access:				oMimicBanRecord.Access,
							Nickname:			string(byNickname),
							Reason:				fmt.Sprintf("Duplicate of banned account (%s)", oMimicBanRecord.SteamID64),
							BanLength:			oMimicBanRecord.BanLength,
							RequestedBy:		oMimicBanRecord.BannedBySteamID64,
						}
						bans.ChanBanManual <- oBanReq;
					}
				}
			}
		} else {
			mapResponse["error"] = "Bad parameters";
		}
	} else {
		mapResponse["error"] = "Bad auth";
	}
}
