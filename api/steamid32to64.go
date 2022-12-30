package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"regexp"
	"fmt"
)


func HttpSteamID32to64(c *gin.Context) {

	var sResponse string;

	sSteamID32 := c.Query("steamid32");

	bSteamID32Valid, _ := regexp.MatchString(`^STEAM_[0-5]:[01]:\d+$`, c.Query("steamid32"));
	if (bSteamID32Valid) {

		arSteamID32 := strings.Split(sSteamID32, ":");
		iY, _ := strconv.ParseInt(arSteamID32[1], 10, 64);
		iZ, _ := strconv.ParseInt(arSteamID32[2], 10, 64);
		sSteamID64 := fmt.Sprintf("%d", 0x0110000100000000 + (iZ * 2) + iY);
		sResponse = sSteamID64;

	} else {
		sResponse = "null";
	}

	c.String(200, sResponse);
}
