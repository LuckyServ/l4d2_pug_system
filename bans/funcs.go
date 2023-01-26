package bans

import (
	"../settings"
	"../database"
	"../players"
	"../utils"
	"../smurf"
	"encoding/base64"
	"time"
	"strings"
	"fmt"
)




func BanRagequitter(oBanReq EntAutoBanReq) { //expensive

	if (!strings.HasPrefix(oBanReq.SteamID64, "7")) { //extra check, just in case
		return;
	}
	
	var iAccess int;
	players.MuPlayers.Lock();
	pPlayer, bFound := players.MapPlayers[oBanReq.SteamID64];
	if (bFound) {
		iAccess = pPlayer.Access;
	}
	players.MuPlayers.Unlock();

	if (iAccess == 4) { //cant ban admin
		return;
	}

	arKnownAccs := smurf.GetKnownAccounts(oBanReq.SteamID64);

	var iCountPrevAutoBans int;
	iSize := len(ArrayBanRecords);
	i64CurTime := time.Now().UnixMilli();
	for i := iSize - 1; i >= 0; i-- {
		if (utils.GetStringIdxInArray(ArrayBanRecords[i].SteamID64, arKnownAccs) != -1) {
			if (i64CurTime - settings.BanHistoryForgetIn >= ArrayBanRecords[i].CreatedAt) {
				break;
			} else {
				if (ArrayBanRecords[i].BannedBySteamID64 == "auto") {
					iCountPrevAutoBans++;
				}
			}
		}
	}

	if (iAccess >= -1) {
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
			Access:				-2,
			BannedBySteamID64:	"auto",
			CreatedAt:			i64BannedAt,
			BanLength:			i64BanLength,
			BanReasonBase64:	sBanReason,
		};

		AddRecord(oBanRecord);

		ApplyBanToPlayer(oBanReq.SteamID64, -2, sBanReason, i64BannedAt, i64BanLength, 0);
		time.Sleep(2 * time.Millisecond);
	}
}

func BanManual(oBanReq EntManualBanReq) {
	if (!strings.HasPrefix(oBanReq.SteamID64, "7")) { //extra check, just in case
		return;
	}

	var iAccess int;

	players.MuPlayers.Lock();
	pPlayer, bFound := players.MapPlayers[oBanReq.SteamID64];
	if (bFound) {
		if (pPlayer.Access > 0) { // moderator cant ban another moderator
			players.MuPlayers.Unlock();
			return;
		}
		iAccess = pPlayer.Access;
	}
	players.MuPlayers.Unlock();

	if (iAccess >= -1) {
		i64BannedAt := time.Now().UnixMilli();
		sBanReason := base64.StdEncoding.EncodeToString([]byte(oBanReq.Reason));
		sNickname := base64.StdEncoding.EncodeToString([]byte(oBanReq.Nickname));
		oBanRecord := EntBanRecord{
			NicknameBase64:		sNickname,
			SteamID64:			oBanReq.SteamID64,
			Access:				oBanReq.Access,
			BannedBySteamID64:	oBanReq.RequestedBy,
			CreatedAt:			i64BannedAt,
			BanLength:			oBanReq.BanLength,
			BanReasonBase64:	sBanReason,
		};

		AddRecord(oBanRecord);

		ApplyBanToPlayer(oBanReq.SteamID64, oBanReq.Access, sBanReason, i64BannedAt, oBanReq.BanLength, 0);
		time.Sleep(2 * time.Millisecond);
	}
}

func ApplyBanToPlayer(sSteamID64 string, iAccess int, sBanReason string, i64BannedAt int64, i64BanLength int64, i64AcceptedAt int64) {
	players.MuPlayers.Lock();
	pPlayer, bFound := players.MapPlayers[sSteamID64];
	if (bFound) {

		pPlayer.BanReason = sBanReason;
		pPlayer.BannedAt = i64BannedAt;
		pPlayer.BanAcceptedAt = i64AcceptedAt;
		pPlayer.BanLength = i64BanLength;
		pPlayer.Access = iAccess;
		go database.UpdatePlayer(database.DatabasePlayer{
			SteamID64:				pPlayer.SteamID64,
			NicknameBase64:			pPlayer.NicknameBase64,
			AvatarSmall:			pPlayer.AvatarSmall,
			AvatarBig:				pPlayer.AvatarBig,
			Mmr:					pPlayer.Mmr,
			MmrUncertainty:			pPlayer.MmrUncertainty,
			LastGameResult:			pPlayer.LastGameResult,
			Access:					pPlayer.Access,
			ProfValidated:			pPlayer.ProfValidated,
			RulesAccepted:			pPlayer.RulesAccepted,
			Twitch:					pPlayer.Twitch,
			CustomMapsConfirmed:	pPlayer.CustomMapsConfirmed,
			LastCampaignsPlayed:	strings.Join(pPlayer.LastCampaignsPlayed, "|"),
			});
		
		players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
	}
	players.MuPlayers.Unlock();
}


func AcceptBan(sSteamID64 string) { //expensive, locks Players

	players.MuPlayers.Lock();

	pPlayer, bFound := players.MapPlayers[sSteamID64];
	if (bFound && pPlayer.Access <= -2 && pPlayer.BanAcceptedAt == 0) {
		pPlayer.BanAcceptedAt = time.Now().UnixMilli();

		iSize := len(ArrayBanRecords);
		for i := iSize - 1; i >= 0; i-- {
			if (ArrayBanRecords[i].CreatedAt == pPlayer.BannedAt) {

				oBanRecord := ArrayBanRecords[i];
				oBanRecord.AcceptedAt = pPlayer.BanAcceptedAt;
				ArrayBanRecords[i] = oBanRecord;

				go database.UpdateBanRecord(database.DatabaseBanRecord{
					NicknameBase64:		oBanRecord.NicknameBase64,
					Access:				oBanRecord.Access,
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
		Access:				oBanRecord.Access,
		BannedBySteamID64:	oBanRecord.BannedBySteamID64,
		CreatedAt:			oBanRecord.CreatedAt,
		AcceptedAt:			oBanRecord.AcceptedAt,
		BanLength:			oBanRecord.BanLength,
		BanReasonBase64:	oBanRecord.BanReasonBase64,
		});
}

func SearchBan(sSteamID64 string) {
	iSize := len(ArrayBanRecords);
	i64CurTime := time.Now().UnixMilli();
	for i := iSize - 1; i >= 0; i-- {
		if (ArrayBanRecords[i].SteamID64 == sSteamID64) {
			if (ArrayBanRecords[i].AcceptedAt == 0 || ArrayBanRecords[i].AcceptedAt + ArrayBanRecords[i].BanLength > i64CurTime) {
				ApplyBanToPlayer(ArrayBanRecords[i].SteamID64, ArrayBanRecords[i].Access, ArrayBanRecords[i].BanReasonBase64, ArrayBanRecords[i].CreatedAt, ArrayBanRecords[i].BanLength, ArrayBanRecords[i].AcceptedAt);
				return;
			}
		}
	}
}

func UnbanManual(sSteamID64 string) {

	players.MuPlayers.Lock();
	pPlayer, bFound := players.MapPlayers[sSteamID64];
	if (bFound && pPlayer.Access <= -2) {
		pPlayer.Access = 0;
		pPlayer.BanReason = "";
		pPlayer.BanAcceptedAt = 0;
		pPlayer.BanLength = 0;
		pPlayer.BannedAt = 0;
		players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
		go database.UpdatePlayer(database.DatabasePlayer{
			SteamID64:				pPlayer.SteamID64,
			NicknameBase64:			pPlayer.NicknameBase64,
			AvatarSmall:			pPlayer.AvatarSmall,
			AvatarBig:				pPlayer.AvatarBig,
			Mmr:					pPlayer.Mmr,
			MmrUncertainty:			pPlayer.MmrUncertainty,
			LastGameResult:			pPlayer.LastGameResult,
			Access:					pPlayer.Access,
			ProfValidated:			pPlayer.ProfValidated,
			RulesAccepted:			pPlayer.RulesAccepted,
			Twitch:					pPlayer.Twitch,
			CustomMapsConfirmed:	pPlayer.CustomMapsConfirmed,
			LastCampaignsPlayed:	strings.Join(pPlayer.LastCampaignsPlayed, "|"),
			});
	}
	players.MuPlayers.Unlock();


	iSize := len(ArrayBanRecords);
	i64CurTime := time.Now().UnixMilli();
	for i := iSize - 1; i >= 0; i-- {
		if (ArrayBanRecords[i].SteamID64 == sSteamID64) {
			if (ArrayBanRecords[i].AcceptedAt == 0 || ArrayBanRecords[i].AcceptedAt + ArrayBanRecords[i].BanLength > i64CurTime) {
				ArrayBanRecords[i].AcceptedAt = i64CurTime - 1;
				ArrayBanRecords[i].BanLength = 1;
				go database.UpdateBanRecord(database.DatabaseBanRecord{
					NicknameBase64:		ArrayBanRecords[i].NicknameBase64,
					Access:				ArrayBanRecords[i].Access,
					SteamID64:			ArrayBanRecords[i].SteamID64,
					BannedBySteamID64:	ArrayBanRecords[i].BannedBySteamID64,
					CreatedAt:			ArrayBanRecords[i].CreatedAt,
					AcceptedAt:			ArrayBanRecords[i].AcceptedAt,
					BanLength:			ArrayBanRecords[i].BanLength,
					BanReasonBase64:	ArrayBanRecords[i].BanReasonBase64,
					});
				return;
			}
		}
	}
}

func DeleteBan(i64CreatedAt int64) {

	iSize := len(ArrayBanRecords);
	for i := iSize - 1; i >= 0; i-- {
		if (ArrayBanRecords[i].CreatedAt == i64CreatedAt) {
			database.DeleteBanRecord(i64CreatedAt);
			ApplyBanToPlayer(ArrayBanRecords[i].SteamID64, 0, "", 0, 0, 0);
			ArrayBanRecords = append(ArrayBanRecords[:i], ArrayBanRecords[i+1:]...);
			return;
		}
	}
}

func RestoreBans() bool {
	arDatabaseBanRecords := database.RestoreBans();
	i64CurTime := time.Now().UnixMilli();
	for _, oDBBanRecord := range arDatabaseBanRecords {
		oBanRecord := EntBanRecord{
			NicknameBase64:		oDBBanRecord.NicknameBase64,
			SteamID64:			oDBBanRecord.SteamID64,
			Access:				oDBBanRecord.Access,
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
				pPlayer.Access = oBanRecord.Access;
			}
		}
	}
	players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
	return true;
}

func BanExcessiveSmurfs(arAccounts []string) {
	iLenAccounts := len(arAccounts);
	if (iLenAccounts > 2) {

		var arToBan, arToBanNames []string;
		players.MuPlayers.RLock();
		for i := 2; i < iLenAccounts; i++ {
			pPlayer, bFound := players.MapPlayers[arAccounts[i]];
			if (bFound && (pPlayer.IsOnline || pPlayer.IsInGame || pPlayer.IsInQueue) && pPlayer.Access >= -1) {
				arToBan = append(arToBan, arAccounts[i]);
				arToBanNames = append(arToBanNames, pPlayer.NicknameBase64);
			}
		}
		players.MuPlayers.RUnlock();


		for i, _ := range arToBan {
			byNickname, _ := base64.StdEncoding.DecodeString(arToBanNames[i]);
			oBanReq := EntManualBanReq{
				SteamID64:			arToBan[i],
				Access:				-2,
				Nickname:			string(byNickname),
				Reason:				fmt.Sprintf("Excessive smurfing, only allowed to play from: %s,%s", arAccounts[0], arAccounts[1]),
				BanLength:			14191200000000,
				RequestedBy:		"smurf",
			}
			BanManual(oBanReq);
		}
	}
}

func BanIfSmurfBanned(arAccounts []string) {
	if (len(arAccounts) > 1) {
		//check if any one of them is banned
		var oMimicBanRecord EntBanRecord;
		i64CurTime := time.Now().UnixMilli();

		iBanlistSize := len(ArrayBanRecords);
		for _, sSteamID64 := range arAccounts {
			for i := iBanlistSize - 1; i >= 0; i-- {
				if (ArrayBanRecords[i].SteamID64 == sSteamID64 && ArrayBanRecords[i].BannedBySteamID64 != "smurf") {
					if (ArrayBanRecords[i].AcceptedAt == 0 || ArrayBanRecords[i].AcceptedAt + ArrayBanRecords[i].BanLength > i64CurTime) {
						oMimicBanRecord = ArrayBanRecords[i];
						break;
					}
				}
			}
			if (oMimicBanRecord.CreatedAt > 0) {
				break;
			}
		}

		if (oMimicBanRecord.CreatedAt > 0) {
			//check if any one of them is online || in game || in queue
			var arToBan []string;
			players.MuPlayers.RLock();
			for _, sSteamID64 := range arAccounts {
				pPlayer, bFound := players.MapPlayers[sSteamID64];
				if (bFound && (pPlayer.IsOnline || pPlayer.IsInGame || pPlayer.IsInQueue) && pPlayer.Access >= -1) {
					arToBan = append(arToBan, sSteamID64);
				}
			}
			players.MuPlayers.RUnlock();


			//ban those found on previous step
			for _, sSteamID64 := range arToBan {
				byNickname, _ := base64.StdEncoding.DecodeString(oMimicBanRecord.NicknameBase64);
				oBanReq := EntManualBanReq{
					SteamID64:			sSteamID64,
					Access:				oMimicBanRecord.Access,
					Nickname:			string(byNickname),
					Reason:				fmt.Sprintf("Duplicate of banned account (%s)", oMimicBanRecord.SteamID64),
					BanLength:			oMimicBanRecord.BanLength,
					RequestedBy:		"bannedsmurf",
				}
				BanManual(oBanReq);
			}
		}
	}
}