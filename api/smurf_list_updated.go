package api

import (
	"github.com/gin-gonic/gin"
	"strings"
	"../bans"
	"../players/auth"
)



func HttpReqSMURFListUpdated(c *gin.Context) {

	mapResponse := make(map[string]interface{});

	mapResponse["success"] = false;
	sAuthKey := c.Query("auth_key");
	if (auth.Backend(sAuthKey)) {
		sAccounts := c.Query("accounts");
		if (sAccounts != "") {
			mapResponse["success"] = true;
			arAccounts := strings.Split(sAccounts, ",");
			bans.ChanAutoBanSmurfs <- arAccounts;
		} else {
			mapResponse["error"] = "Bad parameters";
		}
	} else {
		mapResponse["error"] = "Bad auth";
	}
}
