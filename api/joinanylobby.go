package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../settings"
	"../lobby"
	"time"
	"fmt"
)


func HttpReqJoinAnyLobby(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {

			players.MuPlayers.Lock();

			i64CurTime := time.Now().UnixMilli();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.IsInLobby) {
				mapResponse["error"] = "You are already in a lobby";
			} else if (pPlayer.IsInGame) {
				mapResponse["error"] = "You cant join lobbies, finish your game first";
			} else if (pPlayer.LastLobbyActivity + settings.JoinLobbyCooldown > i64CurTime) {
				mapResponse["error"] = fmt.Sprintf("You cant join lobbies that often. Please wait %d seconds.", ((pPlayer.LastLobbyActivity + settings.JoinLobbyCooldown) - i64CurTime) / 1000);
			} else if (!pPlayer.IsOnline) {
				mapResponse["error"] = "Somehow you are not Online, try to refresh the page";
			} else if (!pPlayer.ProfValidated) {
				mapResponse["error"] = "Please validate your profile first";
			} else if (!pPlayer.RulesAccepted) {
				mapResponse["error"] = "Please accept our rules first";
			} else if (pPlayer.Access <= -2) {
				mapResponse["error"] = "Sorry, you are banned, you gotta wait until it expires";
			} else {
				lobby.MuLobbies.Lock();

				if (lobby.JoinAny(pPlayer)) {
					mapResponse["success"] = true;
					pPlayer.IsAutoSearching = true;
					pPlayer.AutoSearchingSince = time.Now().UnixMilli();
				} else {
					mapResponse["error"] = "Race condition. Try again.";
				}

				lobby.MuLobbies.Unlock();
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
