package bans

import (
	"time"
	"../players"
	"../database"
)


type EntBanRecord struct {
	NicknameBase64		string
	SteamID64			string
	BannedBySteamID64	string
	CreatedAt			int64 //unix timestamp in milliseconds //must keep it unique (which means 1 new ban record per 1ms at most)
	AcceptedAt			int64 //unix timestamp in milliseconds
	BanLength			int64 //unix time in milliseconds
	BanReasonBase64		string
}

type EntAutoBanReq struct {
	SteamID64			string
	NicknameBase64		string
}

var ArrayBanRecords []EntBanRecord;
var ChanUpdateRecord = make(chan EntBanRecord);
var ChanBanRQ = make(chan EntAutoBanReq); //locks Players
var ChanUnban = make(chan bool); //locks Players


func Watchers() {
	go WatchChannels();
	go WatchUnbans();
}

func WatchUnbans() {
	for {
		time.Sleep(60 * time.Second);

		players.MuPlayers.Lock();
		i64CurTime := time.Now().UnixMilli();
		for _, pPlayer := range players.ArrayPlayers {
			if (pPlayer.Access == -2 && pPlayer.BanAcceptedAt > 0 && pPlayer.BanAcceptedAt + pPlayer.BanLength <= i64CurTime) {
				pPlayer.Access = 0;
				pPlayer.BanReason = "";
				pPlayer.BanAcceptedAt = 0;
				pPlayer.BanLength = 0;

				go database.UpdatePlayer(database.DatabasePlayer{
					SteamID64:			pPlayer.SteamID64,
					NicknameBase64:		pPlayer.NicknameBase64,
					Mmr:				pPlayer.Mmr,
					MmrUncertainty:		pPlayer.MmrUncertainty,
					Access:				pPlayer.Access,
					ProfValidated:		pPlayer.ProfValidated,
					RulesAccepted:		pPlayer.RulesAccepted,
					});

				players.I64LastPlayerlistUpdate = i64CurTime;
			}
		}
		players.MuPlayers.Unlock();
	}
}

func WatchChannels() {
	for {
		select {
		case oBanRecord := <-ChanUpdateRecord:
			UpdateRecord(oBanRecord);
		case oBanReq := <-ChanBanRQ: //locks Players
			BanRagequitter(oBanReq);
			time.Sleep(2 * time.Millisecond);
		}
	}
}
