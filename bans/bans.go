package bans

import (
	"time"
	"../players"
	"../database"
	"strings"
)


type EntBanRecord struct {
	NicknameBase64		string
	SteamID64			string
	Access				int
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

type EntManualBanReq struct {
	SteamID64			string
	Access				int
	Nickname			string
	Reason				string
	BanLength			int64
	RequestedBy			string
}

var ArrayBanRecords []EntBanRecord;
var ChanBanRQ = make(chan EntAutoBanReq); //locks Players
var ChanBanManual = make(chan EntManualBanReq); //locks Players
var ChanAcceptBan = make(chan string); //locks Players
var ChanSearchBan = make(chan string); //locks Players
var ChanUnbanManual = make(chan string); //locks Players
var ChanDeleteBan = make(chan int64); //locks Players
var ChanAutoBanSmurfs = make(chan []string); //locks Players
var ChanLock = make(chan bool);
var ChanUnlock = make(chan bool);


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
			if (pPlayer.Access <= -2 && pPlayer.BanAcceptedAt > 0 && pPlayer.BanAcceptedAt + pPlayer.BanLength <= i64CurTime) {
				pPlayer.Access = 0;
				pPlayer.BanReason = "";
				pPlayer.BanAcceptedAt = 0;
				pPlayer.BanLength = 0;
				pPlayer.BannedAt = 0;

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

				players.I64LastPlayerlistUpdate = i64CurTime;
			}
		}
		players.MuPlayers.Unlock();
	}
}

func WatchChannels() {
	for {
		select {
		case oBanReq := <-ChanBanRQ: //locks Players
			BanRagequitter(oBanReq);
		case oBanReq := <-ChanBanManual: //locks Players
			BanManual(oBanReq);
		case sSteamID64 := <-ChanAcceptBan: //locks Players
			AcceptBan(sSteamID64);
		case sSteamID64 := <-ChanSearchBan: //locks Players
			SearchBan(sSteamID64);
		case sSteamID64 := <-ChanUnbanManual: //locks Players
			UnbanManual(sSteamID64);
		case i64CreatedAt := <-ChanDeleteBan:
			DeleteBan(i64CreatedAt);
		case arAccounts := <-ChanAutoBanSmurfs: //locks Players
			BanExcessiveSmurfs(arAccounts);
			BanIfSmurfBanned(arAccounts);
		case <-ChanLock:
			<-ChanUnlock;
		}
	}
}
