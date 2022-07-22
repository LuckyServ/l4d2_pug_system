package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../lobby"
	"../settings"
	"time"
)


func HttpReqJoinLobby(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	sLobbyID := c.Query("lobby_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.IsInLobby) {
				mapResponse["error"] = "You are already in a lobby";
			} else if (pPlayer.LastLobbyActivity + settings.JoinLobbyCooldown > time.Now().UnixMilli()) {
				mapResponse["error"] = "You cant join lobbies that often. Please wait 30 seconds.";
			} else if (!pPlayer.IsOnline) {
				mapResponse["error"] = "Somehow you are not Online, try to refresh the page";
			} else if (!pPlayer.ProfValidated) {
				mapResponse["error"] = "Please validate your profile first";
			} else if (!pPlayer.RulesAccepted) {
				mapResponse["error"] = "Please accept our rules first";
			} else if (pPlayer.Access == -2) {
				mapResponse["error"] = "Sorry, you are banned, you gotta wait until it expires";
			} else if (sLobbyID == "") {
				mapResponse["error"] = "lobby_id parameter isnt set";
			} else {
				//Join lobby
				lobby.MuLobbies.Lock();

				pLobby, bExists := lobby.MapLobbies[sLobbyID];
				if (!bExists) {
					mapResponse["error"] = "This lobby doesnt exist anymore";
				} else if (pLobby.PlayerCount >= 8/*hardcoded for 4v4*/) {
					mapResponse["error"] = "No slots";
				} else if (pPlayer.Mmr < pLobby.MmrMin || pPlayer.Mmr > pLobby.MmrMax) {
					mapResponse["error"] = "Your mmr isnt applicable for this lobby";
				} else {
					if (lobby.Join(pPlayer, sLobbyID)) {
						pPlayer.LastLobbyActivity = time.Now().UnixMilli();
						mapResponse["success"] = true;
					} else {
						mapResponse["error"] = "Race condition. Try again.";
					}
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
	
	c.Header("Access-Control-Allow-Origin", "*");
	c.JSON(200, mapResponse);
}
