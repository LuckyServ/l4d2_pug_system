package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../lobby"
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
				mapResponse["error"] = 2; //already in lobby
			} else if (!pPlayer.IsOnline) {
				mapResponse["error"] = 3; //not online, wtf bro?
			} else if (!pPlayer.ProfValidated) {
				mapResponse["error"] = 4; //profile not validated
			} else if (!pPlayer.RulesAccepted) {
				mapResponse["error"] = 5; //rules not accepted
			} else if (pPlayer.Access == -2) {
				mapResponse["error"] = 6; //banned
			} else if (sLobbyID == "") {
				mapResponse["error"] = 7; //lobby id not set
			} else {
				//Join lobby
				lobby.MuLobbies.Lock();

				pLobby, bExists := lobby.MapLobbies[sLobbyID];
				if (!bExists) {
					mapResponse["error"] = 8; //lobby doesn't exist
				} else if (pLobby.PlayerCount >= 8/*hardcoded for 4v4*/) {
					mapResponse["error"] = 9; //no slots
				} else if (pPlayer.Mmr < pLobby.MmrMin || pPlayer.Mmr > pLobby.MmrMax) {
					mapResponse["error"] = 10; //not applicable mmr
				} else {
					if (lobby.Join(pPlayer, sLobbyID)) {
						mapResponse["success"] = true;
					} else {
						mapResponse["error"] = 11; //? repeat the request
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
