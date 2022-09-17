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
		} else if (iCountPrevAutoBans == 1) {
			i64BanLength = settings.BanRQSecond;
		} else {
			i64BanLength = settings.BanRQThird;
		}

		i64BannedAt := time.Now().UnixMilli();

		sBanReason := base64.StdEncoding.EncodeToString([]byte(settings.BanRQReason));
		oBanRecord := EntBanRecord{
			NicknameBase64:		oBanReq.NicknameBase64,
			SteamID64:			oBanReq.SteamID64,
			BannedBySteamID64:	"auto",
			CreatedAt:			i64CurTime,
			BanLength:			i64BanLength,
			BanReasonBase64:	sBanReason,
		};

		AddRecord(oBanRecord);

		players.MuPlayers.Lock();
		pPlayer, bFound := players.MapPlayers[oBanReq.SteamID64];
		if (bFound) {

			pPlayer.BanReason = sBanReason;
			pPlayer.BannedAt = i64BannedAt;
			pPlayer.BanAcceptedAt = 0;
			pPlayer.BanLength = i64BanLength;
			pPlayer.Access = -2;
			go database.UpdatePlayer(database.DatabasePlayer{
				SteamID64:			pPlayer.SteamID64,
				NicknameBase64:		pPlayer.NicknameBase64,
				Mmr:				pPlayer.Mmr,
				MmrUncertainty:		pPlayer.MmrUncertainty,
				Access:				pPlayer.Access,
				ProfValidated:		pPlayer.ProfValidated,
				RulesAccepted:		pPlayer.RulesAccepted,
				});
			
			players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
		}
		players.MuPlayers.Unlock();
	}
}


func AcceptBan(sSteamID64 string) { //expensive, locks Players

	players.MuPlayers.Lock();

	pPlayer, bFound := players.MapPlayers[sSteamID64];
	if (bFound && pPlayer.Access == -2 && pPlayer.BanAcceptedAt == 0) {
		pPlayer.BanAcceptedAt = time.Now().UnixMilli();

		iSize := len(ArrayBanRecords);
		for i := iSize - 1; i >= 0; i-- {
			if (ArrayBanRecords[i].CreatedAt == pPlayer.BannedAt) {

				oBanRecord := ArrayBanRecords[i];
				oBanRecord.AcceptedAt = pPlayer.BanAcceptedAt;
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

				players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
				players.MuPlayers.Unlock();
				return;
			}
		}
	}

	players.MuPlayers.Unlock();

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

func RestoreBans() bool {
	arDatabaseBanRecords := database.RestoreBans();
	/*if (len(arDatabasePlayers) == 0) {
		return false;
	}*/
	i64CurTime := time.Now().UnixMilli();
	for _, oDBBanRecord := range arDatabaseBanRecords {
		oBanRecord := EntBanRecord{
			NicknameBase64:		oDBBanRecord.NicknameBase64,
			SteamID64:			oDBBanRecord.SteamID64,
			BannedBySteamID64:	oDBBanRecord.BannedBySteamID64,
			CreatedAt:			oDBBanRecord.CreatedAt,
			AcceptedAt:			oDBBanRecord.AcceptedAt,
			BanLength:			oDBBanRecord.BanLength,
			BanReasonBase64:	oDBBanRecord.BanReasonBase64,
		};
		ArrayBanRecords = append(ArrayBanRecords, oBanRecord);

		if (oBanRecord.AcceptedAt == 0 || oBanRecord.AcceptedAt + oBanRecord.BanLength > i64CurTime) {
			pPlayer, bFound := players.MapPlayers[oBanRecord.SteamID64];
			if (bFound) {
				pPlayer.BanReason = oBanRecord.BanReasonBase64;
				pPlayer.BannedAt = oBanRecord.CreatedAt;
				pPlayer.BanAcceptedAt = oBanRecord.AcceptedAt;
				pPlayer.BanLength = oBanRecord.BanLength;
				pPlayer.Access = -2;
			}
		}
	}
	players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
	return true;
}