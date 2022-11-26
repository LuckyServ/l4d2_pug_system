package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../database"
	"time"
)


func HttpReqTicketList(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			i64CurTime := time.Now().UnixMilli();
			if (pPlayer.LastTicketActivity + 1000/*1s*/ > i64CurTime) {
				mapResponse["error"] = "Too often requests, refresh in 1 second";
				players.MuPlayers.Unlock();
			} else if (pPlayer.Access == -3) {
				mapResponse["error"] = "Sorry, you are banned without ability to protest";
				players.MuPlayers.Unlock();
			} else if (pPlayer.Access == 0 && (!pPlayer.ProfValidated || !pPlayer.RulesAccepted)) {
				mapResponse["error"] = "You cant see/create tickets now";
				players.MuPlayers.Unlock();
			} else {
				pPlayer.LastTicketActivity = i64CurTime;
				iAccess := pPlayer.Access;
				players.MuPlayers.Unlock();

				mapResponse["success"] = true;
				mapResponse["opened"] = database.GetOpenedTicketsOfPlayer(oSession.SteamID64);
				mapResponse["closed"] = database.GetClosedTicketsOfPlayer(oSession.SteamID64);
				if (iAccess > 0) {
					mapResponse["as_admin"] = database.GetAdminTickets(iAccess);
				} else {
					mapResponse["as_admin"] = nil;
				}

			}
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
