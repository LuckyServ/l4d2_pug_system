package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../lobby"
	"time"
)


func HttpReqJoinAnyLobby(c *gin.Context) {

	mapResponse := make(map[string]interface{});
	
	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID);
		if (bAuthorized) {

			players.MuPlayers.Lock();

			pPlayer := players.MapPlayers[oSession.SteamID64];
			if (pPlayer.IsInLobby) {
				mapResponse["error"] = "You are already in a lobby";
			} else if (pPlayer.LastLobbyJoin + 30000/*30sec*/ > time.Now().UnixMilli()) {
				mapResponse["error"] = "You cant join lobbies that often. Please wait 30 seconds.";
			} else if (!pPlayer.IsOnline) {
				mapResponse["error"] = "Somehow you are not Online, try to refresh the page";
			} else if (!pPlayer.ProfValidated) {
				mapResponse["error"] = "Please validate your profile first";
			} else if (!pPlayer.RulesAccepted) {
				mapResponse["error"] = "Please accept our rules first";
			} else if (pPlayer.Access == -2) {
				mapResponse["error"] = "Sorry, you are banned, you gotta wait until it expires";
			} else {
				lobby.MuLobbies.Lock();

				arLobbies := lobby.GetJoinableLobbies(pPlayer.Mmr);
				iSize := len(arLobbies);
				if (iSize == 0) {

					if (lobby.Create(pPlayer)) {
						pPlayer.LastLobbyJoin = time.Now().UnixMilli();
						pPlayer.IsAutoSearching = true;
						mapResponse["success"] = true;
					} else {
						mapResponse["error"] = "Race condition on lobby creation. Try again.";
					}

				} else {
					//sort
					if (iSize > 1) {
						bSorted := false;
						for !bSorted {
							bSorted = true;
							for i := 1; i < iSize; i++ {
								if (arLobbies[i].CreatedAt < arLobbies[i - 1].CreatedAt) {
									arLobbies[i], arLobbies[i - 1] = arLobbies[i - 1], arLobbies[i]; //switch
									bSorted = false;
								}
							}
						}
					}
					sLobbyID := arLobbies[0].ID;

					if (lobby.Join(pPlayer, sLobbyID)) {
						pPlayer.LastLobbyJoin = time.Now().UnixMilli();
						pPlayer.IsAutoSearching = true;
						mapResponse["success"] = true;
					} else {
						mapResponse["error"] = "Race condition on lobby join. Try again.";
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
