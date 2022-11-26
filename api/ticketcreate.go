package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../database"
	"../utils"
	"encoding/base64"
	"time"
	"strconv"
	"unicode/utf8"
)


func HttpReqTicketCreate(c *gin.Context) {

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	sTicketText := c.PostForm("ticket_text");
	iTicketType, _ := strconv.Atoi(c.PostForm("ticket_type"));
	sRedirectTo := c.PostForm("redirect_to");
	iTextLen := utf8.RuneCountInString(sTicketText);
	if (errCookieSessID == nil && sCookieSessID != "" && iTextLen > 0 && iTextLen < 10000 && iTicketType > 0 && iTicketType < database.TicketTypeDummy2) {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.PostForm("csrf"));
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			i64CurTime := time.Now().UnixMilli();
			if (pPlayer.LastTicketActivity + 5000/*5s*/ > i64CurTime) {
				players.MuPlayers.Unlock();
			} else if (pPlayer.Access == -3) {
				players.MuPlayers.Unlock();
			} else if (pPlayer.Access == 0 && (!pPlayer.ProfValidated || !pPlayer.RulesAccepted)) {
				players.MuPlayers.Unlock();
			} else {
				pPlayer.LastTicketActivity = i64CurTime;
				players.MuPlayers.Unlock();

				sTicketID := <-utils.ChanUniqueString;

				database.CreateTicket(database.DatabaseTicket{
					TicketID:		sTicketID,
					TicketType:		iTicketType,
					CreatedBy:		oSession.SteamID64,
					CreatedAt:		i64CurTime,
					IsClosed:		false,
				});

				database.CreateMessage(database.DatabaseTicketMessage{
					TicketID:		sTicketID,
					MessageBy:		oSession.SteamID64,
					MessageAt:		i64CurTime,
					MessageBase64:	base64.StdEncoding.EncodeToString([]byte(sTicketText)),
				});

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
