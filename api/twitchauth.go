package api

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"../settings"
)


func HttpTwitchAuth(c *gin.Context) {
	//?provider=twitch
	sHomepage := c.Query("home_page");
	if (sHomepage == "") {
		sHomepage = "https://"+settings.HomeDomain;
	}
	c.SetCookie("home_page", sHomepage, 3600, "/", "", true, false);
	gothic.BeginAuthHandler(c.Writer, c.Request);
}
