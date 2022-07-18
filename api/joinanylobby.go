package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../lobby"
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
				mapResponse["error"] = 2; //already in lobby
			} else if (!pPlayer.IsOnline) {
				mapResponse["error"] = 3; //not online, wtf bro?
			} else if (!pPlayer.ProfValidated) {
				mapResponse["error"] = 4; //profile not validated
			} else if (!pPlayer.RulesAccepted) {
				mapResponse["error"] = 5; //rules not accepted
			} else if (pPlayer.Access == -2) {
				mapResponse["error"] = 6; //banned
			} else {
				lobby.MuLobbies.Lock();

				arLobbies := lobby.GetJoinableLobbies(pPlayer.Mmr);
				iSize := len(arLobbies);
				if (iSize == 0) {

					if (lobby.Create(pPlayer)) {
						mapResponse["success"] = true;
					} else {
						mapResponse["error"] = 7; //Error creating lobby, shouldn't ever happen
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
						mapResponse["success"] = true;
					} else {
						mapResponse["error"] = 8; //Error joining existing lobby
					}
				}
				lobby.MuLobbies.Unlock();
			}
			players.MuPlayers.Unlock();

		} else {
			mapResponse["error"] = 1; //unauthorized
		}
	} else {
		mapResponse["error"] = 1; //unauthorized
	}
	
	c.Header("Access-Control-Allow-Origin", "*");
	c.JSON(200, mapResponse);
}
