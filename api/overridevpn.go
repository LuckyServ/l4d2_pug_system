package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"../smurf"
	"time"
)


func HttpReqOverrideVPN(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	sIP := c.Query("ip");
	sVPN := c.Query("isvpn");

	mapResponse["success"] = false;
	if (sIP != "" && sVPN != "") {
		if (errCookieSessID == nil && sCookieSessID != "") {
			oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
			if (bAuthorized) {
				players.MuPlayers.RLock();
				iAccess := players.MapPlayers[oSession.SteamID64].Access;
				players.MuPlayers.RUnlock();
				if (iAccess == 4) { //only admin can override vpn info

					bVPN := (sVPN == "true");

					smurf.MuVPN.Lock();
					oNewVPNInfo := smurf.EntVPNInfo{
						IsVPN:		bVPN,
						IsInCheck:	false,
						UpdatedAt:	time.Now().Unix(),
					};
					smurf.MapVPNs[sIP] = oNewVPNInfo;
					smurf.MuVPN.Unlock();

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
