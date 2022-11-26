package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../database"
	"time"
	"regexp"
)

func HttpReqTicketMessages(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	sTicketID := c.Query("ticket_id");
	bTicketIDValid, _ := regexp.MatchString(`^[0-9a-z]{1,100}$`, sTicketID);

	mapResponse["success"] = false;
	if (bTicketIDValid) {
		if (errCookieSessID == nil && sCookieSessID != "") {
			oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
			if (bAuthorized) {
				players.MuPlayers.Lock();
				pPlayer := players.MapPlayers[oSession.SteamID64];
				i64CurTime := time.Now().UnixMilli();
				if (pPlayer.LastTicketActivity + 1000/*1s*/ > i64CurTime) {
					mapResponse["error"] = "Too often requests, try again in 1 second";
					players.MuPlayers.Unlock();
				} else if (pPlayer.Access == -3) {
					mapResponse["error"] = "Sorry, you are banned without ability to protest";
					players.MuPlayers.Unlock();
				} else if (pPlayer.Access == 0 && (!pPlayer.ProfValidated || !pPlayer.RulesAccepted)) {
					mapResponse["error"] = "You cant see/create tickets now";
					players.MuPlayers.Unlock();
				} else {
					pPlayer.LastTicketActivity = i64CurTime;
					players.MuPlayers.Unlock();

					mapResponse["success"] = true;
					mapResponse["ticket_id"] = sTicketID;
					mapResponse["messages"] = database.GetMessagesOfTicket(sTicketID);

				}
			} else {
				mapResponse["error"] = "Please authorize first";
			}
		} else {
			mapResponse["error"] = "Please authorize first";
		}
	} else {
		mapResponse["error"] = "Bad ticket_id parameter";
	}
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
