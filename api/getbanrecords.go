package api

import (
	"github.com/gin-gonic/gin"
	"../settings"
	"../bans"
	"../utils"
	"strconv"
	"encoding/base64"
)

type BanRecordResponse struct {
	NicknameBase64		string		`json:"base64name"`
	SteamID64			string		`json:"steamid64"`
	CreatedAt			int64		`json:"created_at"`
	BannedBySteamID64	string		`json:"banned_by"`
	AcceptedAt			int64		`json:"accepted_at"`
	BanLength			int64		`json:"ban_length"`
	BanReasonBase64		string		`json:"base64reason"`
}


func HttpReqGetBanRecords(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	iPage, _ := strconv.Atoi(c.Query("page"));
	sSearch := c.Query("search");

	mapResponse["success"] = false;
	mapResponse["page"] = iPage;

	var arFilteredBanRecords []bans.EntBanRecord;
	bans.ChanLock <- true;

	if (sSearch == "") {
		arFilteredBanRecords = bans.ArrayBanRecords;
	} else {
		for _, oBanRecord := range bans.ArrayBanRecords {
			byNickname, _ := base64.StdEncoding.DecodeString(oBanRecord.NicknameBase64);
			sNickname := string(byNickname);
			byReason, _ := base64.StdEncoding.DecodeString(oBanRecord.BanReasonBase64);
			sReason := string(byReason);
			if (utils.StringContainsCI(sSearch, sNickname) || utils.StringContainsCI(sSearch, oBanRecord.SteamID64) || utils.StringContainsCI(sSearch, sReason)) {
				arFilteredBanRecords = append(arFilteredBanRecords, oBanRecord);
			}
		}
	}

	iMaxBanRecords := len(arFilteredBanRecords);
	iStartItem := (iMaxBanRecords - 1) - (iPage * settings.BanListPagination);
	if (iStartItem < iMaxBanRecords && iStartItem >= 0) {
		var arBanRecordsResp []BanRecordResponse;

		iEndItem := (iStartItem - settings.BanListPagination) + 1;
		if (iEndItem < 0) {
			iEndItem = 0;
		}
		for i := iStartItem; i >= iEndItem; i-- {
			arBanRecordsResp = append(arBanRecordsResp, BanRecordResponse{
				NicknameBase64:		arFilteredBanRecords[i].NicknameBase64,
				SteamID64:			arFilteredBanRecords[i].SteamID64,
				CreatedAt:			arFilteredBanRecords[i].CreatedAt,
				BannedBySteamID64:	arFilteredBanRecords[i].BannedBySteamID64,
				AcceptedAt:			arFilteredBanRecords[i].AcceptedAt,
				BanLength:			arFilteredBanRecords[i].BanLength,
				BanReasonBase64:	arFilteredBanRecords[i].BanReasonBase64,
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
