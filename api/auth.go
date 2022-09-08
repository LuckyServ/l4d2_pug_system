package api

import (
	"github.com/gin-gonic/gin"
	"../settings"
)


func HttpReqAuth(c *gin.Context) {
	sHomepage := c.Query("home_page");
	if (sHomepage == "") {
		sHomepage = "https://"+settings.HomeDomain;
	}
	c.SetCookie("home_page", sHomepage, 3600, "/", "", true, false);
	c.Redirect(303, "https://steamcommunity.com/openid/login?openid.identity=http://specs.openid.net/auth/2.0/identifier_select&openid.claimed_id=http://specs.openid.net/auth/2.0/identifier_select&openid.ns=http://specs.openid.net/auth/2.0&openid.mode=checkid_setup&openid.realm=https://"+settings.BackendDomain+"&openid.return_to=https://"+settings.BackendDomain+"/openidcallback");
}
