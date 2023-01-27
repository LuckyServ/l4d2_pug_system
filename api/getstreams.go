package api

import (
	"github.com/gin-gonic/gin"
	"../players"
	"../streams"
	//"fmt"
	"time"
	"../players/auth"
)


func HttpReqGetStreams(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");

	mapResponse["authorized"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			mapResponse["authorized"] = true;
			players.MuPlayers.RLock();
			pPlayer := players.MapPlayers[oSession.SteamID64];
			mapResponse["steamid64"] = oSession.SteamID64;
			mapResponse["my_stream"] = pPlayer.Twitch;
			players.MuPlayers.RUnlock();
		}
	}

	streams.MuStreams.RLock();
	i64CurTime := time.Now().UnixMilli();
	arStreams := make([]streams.TwitchStream, len(streams.ArrayStreams));
	copy(arStreams, streams.ArrayStreams);
	streams.MuStreams.RUnlock();


	mapResponse["success"] = true;
	mapResponse["streams"] = arStreams;

	mapResponse["updated_at"] = i64CurTime;

	
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
