package chat

import (
	"../settings"
	"../database"
	"encoding/base64"
)

type EntChatMsg struct {
	TimeStamp		int64 //ms
	Text			string
	SteamID64		string //who sent it
	NicknameBase64	string //who sent it
	AvatarSmall		string //who sent it
}

var ArrayChatMsgs []EntChatMsg;
var ChanSend = make(chan EntChatMsg);
var ChanRead = make(chan []EntChatMsg);
var ChanDelete = make(chan string);
var I64LastGlobalChatUpdate int64;




func ChannelWatchers() {
	for {
		select {
		case oChatMsg := <-ChanSend:

			if (len(ArrayChatMsgs) >= settings.ChatStoreMaxMsgs) {
				ArrayChatMsgs = ArrayChatMsgs[1:];
			}
			ArrayChatMsgs = append(ArrayChatMsgs, oChatMsg);
			go database.PublicChatLog(oChatMsg.TimeStamp, oChatMsg.NicknameBase64, oChatMsg.SteamID64, base64.StdEncoding.EncodeToString([]byte(oChatMsg.Text)));

		case ChanRead <- func()([]EntChatMsg) {
			arChatMsgs := make([]EntChatMsg, len(ArrayChatMsgs));
			copy(arChatMsgs, ArrayChatMsgs);
			return arChatMsgs;
		}():
		case sSteamID64 := <-ChanDelete:
			bNoneLeft := false;
			for !bNoneLeft {
				bNoneLeft = true;
				for i, oChatSmg := range ArrayChatMsgs {
					if (oChatSmg.SteamID64 == sSteamID64) {
						bNoneLeft = false;
						ArrayChatMsgs = append(ArrayChatMsgs[:i], ArrayChatMsgs[i+1:]...);
						break;
					}
				}
			}
		}
	}
}


func AddChatMessage(oChatMsg EntChatMsg) {

	if (len(ArrayChatMsgs) >= settings.ChatStoreMaxMsgs) {
		ArrayChatMsgs = ArrayChatMsgs[1:];
	}

	ArrayChatMsgs = append(ArrayChatMsgs, oChatMsg);
}
