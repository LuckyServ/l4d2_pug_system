package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../database"
	"encoding/base64"
	"time"
	"unicode/utf8"
	"regexp"
)


func HttpReqTicketReply(c *gin.Context) {

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	sMessageText := c.PostForm("message_text");
	sTicketID := c.PostForm("ticket_id");
	bTicketIDValid, _ := regexp.MatchString(`^[0-9a-z]{1,100}$`, sTicketID);
	sRedirectTo := c.PostForm("redirect_to");
	iTextLen := utf8.RuneCountInString(sMessageText);
	if (errCookieSessID == nil && sCookieSessID != "" && iTextLen > 0 && iTextLen < 10000 && bTicketIDValid) {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.PostForm("csrf"));
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			i64CurTime := time.Now().UnixMilli();
			if (pPlayer.LastTicketActivity + 1000/*1s*/ > i64CurTime) {
				players.MuPlayers.Unlock();
			} else if (pPlayer.Access < 1) {
				players.MuPlayers.Unlock();
			} else {
				pPlayer.LastTicketActivity = i64CurTime;
				players.MuPlayers.Unlock();

				database.CreateMessage(database.DatabaseTicketMessage{
					TicketID:		sTicketID,
					MessageBy:		oSession.SteamID64,
					MessageAt:		i64CurTime,
					MessageBase64:	base64.StdEncoding.EncodeToString([]byte(sMessageText)),
				});

				database.CloseTicket(sTicketID);

			}
		}
	}
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	if (sRedirectTo == "") {
		c.String(303, "ok");
	} else {
		c.Redirect(303, sRedirectTo);
	}
}
