package api

import (
	"github.com/gin-gonic/gin"
	"../settings"
	"../bans"
	"strconv"
)

type BanRecordResponse struct {
	NicknameBase64		string		`json:"base64name"`
	SteamID64			string		`json:"steamid64"`
	CreatedAt			int64		`json:"created_at"`
	AcceptedAt			int64		`json:"accepted_at"`
	BanLength			int64		`json:"ban_length"`
	BanReasonBase64		string		`json:"base64reason"`
}


func HttpReqGetBanRecords(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	iPage, _ := strconv.Atoi(c.Query("page"));

	mapResponse["success"] = false;
	mapResponse["page"] = iPage;
	bans.ChanLock <- true;
	iMaxBanRecords := len(bans.ArrayBanRecords);
	iStartItem := (iMaxBanRecords - 1) - (iPage * settings.BanListPagination);
	if (iStartItem < iMaxBanRecords && iStartItem >= 0) {
		var arBanRecordsResp []BanRecordResponse;

		iEndItem := (iStartItem - settings.BanListPagination) + 1;
		if (iEndItem < 0) {
			iEndItem = 0;
		}
		for i := iStartItem; i >= iEndItem; i-- {
			arBanRecordsResp = append(arBanRecordsResp, BanRecordResponse{
				NicknameBase64:		bans.ArrayBanRecords[i].NicknameBase64,
				SteamID64:			bans.ArrayBanRecords[i].SteamID64,
				CreatedAt:			bans.ArrayBanRecords[i].CreatedAt,
				AcceptedAt:			bans.ArrayBanRecords[i].AcceptedAt,
				BanLength:			bans.ArrayBanRecords[i].BanLength,
				BanReasonBase64:	bans.ArrayBanRecords[i].BanReasonBase64,
			});
		}
		bans.ChanUnlock <- true;
		mapResponse["success"] = true;
		mapResponse["bans"] = arBanRecordsResp;
	} else {bans.ChanUnlock <- true;};





	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
