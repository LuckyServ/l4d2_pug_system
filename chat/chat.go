package chat

import (
	"../settings"
	"time"
)

type EntChatMsg struct {
	TimeStamp		int64
	Text			string
	SteamID64		string //who sent it
	NicknameBase64	string //who sent it
	AvatarSmall		string //who sent it
}

var ArrayChatMsgs []EntChatMsg;
var ChanSend = make(chan EntChatMsg);
var ChanGetUniqTime = make(chan int64);
var ChanRead = make(chan []EntChatMsg);
var I64LastGlobalChatUpdate int64;




func ChannelWatchers() {
	for {
		select {
		case oChatMsg := <-ChanSend:

			if (len(ArrayChatMsgs) >= settings.ChatStoreMaxMsgs) {
				ArrayChatMsgs = ArrayChatMsgs[1:];
			}
			ArrayChatMsgs = append(ArrayChatMsgs, oChatMsg);

		case ChanRead <- func()([]EntChatMsg) {
			arChatMsgs := make([]EntChatMsg, len(ArrayChatMsgs));
			copy(arChatMsgs, ArrayChatMsgs);
			return arChatMsgs;
		}():
		case ChanGetUniqTime <- func()(int64) {
			time.Sleep(1 * time.Millisecond);
			return time.Now().UnixMilli();
		}():
		}
	}
}


func AddChatMessage(oChatMsg EntChatMsg) {

	if (len(ArrayChatMsgs) >= settings.ChatStoreMaxMsgs) {
		ArrayChatMsgs = ArrayChatMsgs[1:];
	}

	ArrayChatMsgs = append(ArrayChatMsgs, oChatMsg);
}
