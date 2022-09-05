package chat

import (
	"../settings"
)

type EntChatMsg struct {
	Text			string
	SteamID64		string //who sent it
	NicknameBase64	string //who sent it
}

var ArrayChatMsgs []EntChatMsg;
var ChanSend = make(chan EntChatMsg);
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
		}
	}
}


func AddChatMessage(oChatMsg EntChatMsg) {

	if (len(ArrayChatMsgs) >= settings.ChatStoreMaxMsgs) {
		ArrayChatMsgs = ArrayChatMsgs[1:];
	}

	ArrayChatMsgs = append(ArrayChatMsgs, oChatMsg);
}
