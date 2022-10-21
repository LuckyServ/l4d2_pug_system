package api

import (
	"github.com/gin-gonic/gin"
	"../chat"
	"encoding/base64"
	"time"
	"fmt"
)


type ChatMsgResp struct {
	TextBase64		string	`json:"base64text"`
	SteamID64		string	`json:"steamid64"`
	NicknameBase64	string	`json:"base64name"`
	AvatarSmall		string	`json:"avatar_small"`
}


func HttpReqGetGlobalChat(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	arChatMsgs := <-chat.ChanRead;

	var arRespChat []ChatMsgResp;
	for _, oMsg := range arChatMsgs {
		arRespChat = append(arRespChat, ChatMsgResp{
			TextBase64:		base64.StdEncoding.EncodeToString([]byte(oMsg.Text)),
			SteamID64:		oMsg.SteamID64,
			NicknameBase64:	oMsg.NicknameBase64,
			AvatarSmall:	oMsg.AvatarSmall,
		});
	}

	mapResponse["success"] = true;
	mapResponse["messages"] = arRespChat;

	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.SetCookie("globalchat_updated_at", fmt.Sprintf("%d", time.Now().UnixMilli()), 2592000, "/", "", true, false);
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
