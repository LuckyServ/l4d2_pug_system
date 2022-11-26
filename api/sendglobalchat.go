package api

import (
	"github.com/gin-gonic/gin"
	"../settings"
	"../chat"
	"../players"
	"../bans"
	"../smurf"
	"../players/auth"
	"unicode/utf8"
	"time"
)


func HttpReqSendGlobalChat(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sChatMsg := c.Query("text");

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			i64CurTime := time.Now().UnixMilli();
			if (pPlayer.Access <= -1) {
				mapResponse["error"] = "Sorry, you are banned, you have to wait until it expires";
			} else if (pPlayer.LastChatMessage + settings.ChatMsgDelay > i64CurTime) {
				mapResponse["error"] = "no_alert";
			} else if (chat.I64LastGlobalChatUpdate + 5/*once per 5ms*/ > i64CurTime) {
				mapResponse["error"] = "no_alert";
			} else if (!pPlayer.ProfValidated) {
				mapResponse["error"] = "Please validate your profile first";
			} else if (!pPlayer.RulesAccepted) {
				mapResponse["error"] = "Please accept our rules first";
			} else if (!pPlayer.IsOnline) {
				mapResponse["error"] = "Somehow you are not Online, try to refresh the page";
			} else {
				iTextLen := utf8.RuneCountInString(sChatMsg);
				if (iTextLen > 0 && iTextLen <= settings.ChatMaxChars) {
					mapResponse["success"] = true;
					oMessage := chat.EntChatMsg{
						TimeStamp:		<-chat.ChanGetUniqTime,
						Text:			sChatMsg,
						SteamID64:		pPlayer.SteamID64,
						NicknameBase64:	pPlayer.NicknameBase64,
						AvatarSmall:	pPlayer.AvatarSmall,
					};
					chat.ChanSend <- oMessage;
					i64CurTime = time.Now().UnixMilli();
					pPlayer.LastChatMessage = i64CurTime;
					chat.I64LastGlobalChatUpdate = i64CurTime;
					go func(sSteamID64 string)() {bans.ChanAutoBanSmurfs <- smurf.GetKnownAccounts(sSteamID64);}(oSession.SteamID64);
				} else {
					mapResponse["error"] = "Bad message size";
				}
			}
			players.MuPlayers.Unlock();
		} else {
			mapResponse["error"] = "Please authorize first";
		}
	} else {
		mapResponse["error"] = "Please authorize first";
	}


	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
