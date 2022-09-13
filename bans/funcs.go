package bans

import (
	"../settings"
	"../database"
	"../players"
	"encoding/base64"
	"time"
	"strings"
)




func BanRagequitter(oBanReq EntAutoBanReq) { //expensive

	if (!strings.HasPrefix(oBanReq.SteamID64, "7")) { //extra check, just in case
		return;
	}
	
	var iCountPrevAutoBans int;
	var bIsBannedNow bool;

	i64CurTime := time.Now().UnixMilli();
	iSize := len(ArrayBanRecords);
	for i := iSize - 1; i >= 0; i-- {
		if (ArrayBanRecords[i].SteamID64 == oBanReq.SteamID64) {
			if (ArrayBanRecords[i].AcceptedAt == 0 || ArrayBanRecords[i].AcceptedAt + ArrayBanRecords[i].BanLength > i64CurTime) {
				if (!bIsBannedNow) {
						bIsBannedNow = true;
					}
				}
			if (ArrayBanRecords[i].CreatedAt + settings.BanHistoryForgetIn > i64CurTime && ArrayBanRecords[i].BannedBySteamID64 == "auto") {
				iCountPrevAutoBans++;
			}
		}
	}

	if (!bIsBannedNow) {
		var i64BanLength int64;
		if (iCountPrevAutoBans == 0) {
			i64BanLength = settings.BanRQFirst;
		} else {
			i64BanLength = settings.BanRQSecond;
		}

		sBanReason := base64.StdEncoding.EncodeToString([]byte(settings.BanRQReason));
		oBanRecord := EntBanRecord{
			NicknameBase64:		oBanReq.NicknameBase64,
			SteamID64:			oBanReq.SteamID64,
			BannedBySteamID64:	"auto",
			CreatedAt:			time.Now().UnixMilli(),
			BanLength:			i64BanLength,
			BanReasonBase64:	sBanReason,
		};

		AddRecord(oBanRecord);

		players.MuPlayers.Lock();
		pPlayer, bFound := players.MapPlayers[oBanReq.SteamID64];
		if (bFound) {

			pPlayer.BanReason = sBanReason;
			pPlayer.BanAcceptedAt = 0;
			pPlayer.BanLength = i64BanLength;
			pPlayer.Access = -2;
			
		}
		players.MuPlayers.Unlock();
	}
}


func AddRecord(oBanRecord EntBanRecord) {
	ArrayBanRecords = append(ArrayBanRecords, oBanRecord);
	go database.AddBanRecord(database.DatabaseBanRecord{
		NicknameBase64:		oBanRecord.NicknameBase64,
		SteamID64:			oBanRecord.SteamID64,
		BannedBySteamID64:	oBanRecord.BannedBySteamID64,
		CreatedAt:			oBanRecord.CreatedAt,
		AcceptedAt:			oBanRecord.AcceptedAt,
		BanLength:			oBanRecord.BanLength,
		BanReasonBase64:	oBanRecord.BanReasonBase64,
		});
}


func UpdateRecord(oBanRecord EntBanRecord) { //expensive
	iSize := len(ArrayBanRecords);
	for i := iSize - 1; i >= 0; i-- {
		if (ArrayBanRecords[i].CreatedAt == oBanRecord.CreatedAt) {
			ArrayBanRecords[i] = oBanRecord;
			go database.UpdateBanRecord(database.DatabaseBanRecord{
				NicknameBase64:		oBanRecord.NicknameBase64,
				SteamID64:			oBanRecord.SteamID64,
				BannedBySteamID64:	oBanRecord.BannedBySteamID64,
				CreatedAt:			oBanRecord.CreatedAt,
				AcceptedAt:			oBanRecord.AcceptedAt,
				BanLength:			oBanRecord.BanLength,
				BanReasonBase64:	oBanRecord.BanReasonBase64,
				});
			return;
		}
	}
}
