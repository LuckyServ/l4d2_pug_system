package bans

import (
	"time"
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