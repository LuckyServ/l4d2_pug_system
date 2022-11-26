package api

import (
	"github.com/gin-gonic/gin"
	"../players/auth"
	"../players"
	"time"
)

func HttpReqLogout(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	mapResponse["success"] = false;
	if (errCookieSessID == nil && sCookieSessID != "") {
		_, bAuthorized := auth.GetSession(sCookieSessID, c.Query("csrf"));
		if (bAuthorized) {
			if (auth.RemoveSession(sCookieSessID)) {
				mapResponse["success"] = true;
				players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
				c.SetCookie("session_id", "", 1, "/", "", true, false);
			} else {
				mapResponse["error"] = "Error";
			}
		} else {
			mapResponse["error"] = "You are not authorized";
		}
	} else {
		mapResponse["error"] = "You are not authorized";
	}

	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("origin"));
	c.Header("Access-Control-Allow-Credentials", "true");
	c.JSON(200, mapResponse);
}
