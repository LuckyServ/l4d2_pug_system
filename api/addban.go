package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../bans"
	"strconv"
)


func HttpReqAddBan(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	sSteamID64 := c.Query("steamid64");
	sNickname := c.Query("nickname");
	sReason := c.Query("reason");
	iBanType, _ := strconv.Atoi(c.Query("bantype"));
	i64BanLength, _ := strconv.ParseInt(c.Query("banlength"), 10, 64);

	mapResponse["success"] = false;
	if (sSteamID64 != "" && sNickname != "" && i64BanLength > 0 && iBanType <= -2 && iBanType >= -3) {
		if (errCookieSessID == nil && sCookieSessID != "") {
			oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
			if (bAuthorized) {
				players.MuPlayers.RLock();
				iAccess := players.MapPlayers[oSession.SteamID64].Access;
				players.MuPlayers.RUnlock();
				if (iAccess > 0) { //moderator or admin

					bans.ChanBanManual <- bans.EntManualBanReq{
						SteamID64:			sSteamID64,
						Access:				iBanType,
						Nickname:			sNickname,
						Reason:				sReason,
						BanLength:			i64BanLength,
						RequestedBy:		oSession.SteamID64,
					};

					mapResponse["success"] = true;
				} else {
					mapResponse["error"] = "You dont have access to this command";
				}
			} else {
				mapResponse["error"] = "Please authorize first";
			}
		} else {
			mapResponse["error"] = "Please authorize first";
		}
	} else {
		mapResponse["error"] = "Bad parameters";
	}
	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
