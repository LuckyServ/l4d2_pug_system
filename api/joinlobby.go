package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../lobby"
	"../settings"
	"time"
	"fmt"
	"encoding/base64"
	"../smurf"
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
						mapResponse["success"] = true;
					} else {
						mapResponse["error"] = "Race condition. Try again.";
					}
				}

				lobby.MuLobbies.Unlock();

				sCookieUniqueKey, _ := c.Cookie("auth2");
				byNickname, _ := base64.StdEncoding.DecodeString(pPlayer.NicknameBase64);
				go smurf.AnnounceIPAndKey(pPlayer.SteamID64, c.ClientIP(), string(byNickname), sCookieUniqueKey);

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
