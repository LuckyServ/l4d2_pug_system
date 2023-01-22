package api

import (
	"github.com/gin-gonic/gin"
	"../settings"
	//"fmt"
	"time"
	"sync"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitch"
	"../players"
	"../players/auth"
	"../database"
	"strings"
)


var mapIPsTwitch = make(map[string]int);
var MuAuthTwitch sync.RWMutex;

func SetupOpenID() {
	goth.UseProviders(twitch.New(settings.TwitchClientID, settings.TwitchSecret, "https://"+settings.BackendDomain + "/twitchcallback"));
}

//Limit authorizations per hour per IP
func TwitchAuthRatelimits() {
	for {
		time.Sleep(3600 * time.Second); //1 hour
		MuAuthTwitch.Lock();
		mapIPsTwitch = make(map[string]int);
		MuAuthTwitch.Unlock();
	}
}

func TwitchGetAuthCount(sClientIP string) int {
	MuAuthTwitch.RLock();
	iCount, bExists := mapIPsTwitch[sClientIP];
	MuAuthTwitch.RUnlock();
	if (bExists) {
		return iCount;
	} else {
		return 0;
	}
}

func TwitchIncreaseAuthCount(sClientIP string) {
	MuAuthTwitch.Lock();
	iCount, bExists := mapIPsTwitch[sClientIP];
	if (bExists) {
		mapIPsTwitch[sClientIP] = iCount + 1;
	} else {
		mapIPsTwitch[sClientIP] = 1;
	}
	MuAuthTwitch.Unlock();
}


func HttpTwitchOpenIDCallback(c *gin.Context) {

	//Ratelimits
	sClientIP := c.ClientIP();
	iCount := TwitchGetAuthCount(sClientIP);
	if (iCount >= settings.AuthPerHour) {
		c.String(200, "Too many authorization requests. Wait an hour before trying again.");
		return;
	}
	TwitchIncreaseAuthCount(sClientIP);

	sCookieSessID, errCookieSessID := c.Cookie("session_id");
	sHomepage, errHomepage := c.Cookie("home_page");
	if (errHomepage != nil || sHomepage == "") {
		sHomepage = "https://"+settings.HomeDomain;
	}
	defer c.Redirect(303, sHomepage);

	oAuthResponse, errAuthResponce := gothic.CompleteUserAuth(c.Writer, c.Request);
	if (errAuthResponce != nil) {
		return;
	}
	if (oAuthResponse.Provider != "twitch" || oAuthResponse.NickName == "" || oAuthResponse.UserID == "") {
		return;
	}

	if (errCookieSessID == nil && sCookieSessID != "") {
		oSession, bAuthorized := auth.GetSessionNoCSRF(sCookieSessID);
		if (bAuthorized) {
			players.MuPlayers.Lock();
			pPlayer := players.MapPlayers[oSession.SteamID64];

			for _, pPlayerSearch := range players.ArrayPlayers {
				if (pPlayerSearch.Twitch == oAuthResponse.UserID && pPlayerSearch != pPlayer) {
					pPlayerSearch.Twitch = "";
					go database.UpdatePlayer(database.DatabasePlayer{
						SteamID64:				pPlayerSearch.SteamID64,
						NicknameBase64:			pPlayerSearch.NicknameBase64,
						AvatarSmall:			pPlayerSearch.AvatarSmall,
						AvatarBig:				pPlayerSearch.AvatarBig,
						Mmr:					pPlayerSearch.Mmr,
						MmrUncertainty:			pPlayerSearch.MmrUncertainty,
						LastGameResult:			pPlayerSearch.LastGameResult,
						Access:					pPlayerSearch.Access,
						ProfValidated:			pPlayerSearch.ProfValidated,
						RulesAccepted:			pPlayerSearch.RulesAccepted,
						Twitch:					pPlayerSearch.Twitch,
						CustomMapsConfirmed:	pPlayerSearch.CustomMapsConfirmed,
						LastCampaignsPlayed:	strings.Join(pPlayerSearch.LastCampaignsPlayed, "|"),
						});
				}
			}

			pPlayer.Twitch = oAuthResponse.UserID;
			go database.UpdatePlayer(database.DatabasePlayer{
				SteamID64:				pPlayer.SteamID64,
				NicknameBase64:			pPlayer.NicknameBase64,
				AvatarSmall:			pPlayer.AvatarSmall,
				AvatarBig:				pPlayer.AvatarBig,
				Mmr:					pPlayer.Mmr,
				MmrUncertainty:			pPlayer.MmrUncertainty,
				LastGameResult:			pPlayer.LastGameResult,
				Access:					pPlayer.Access,
				ProfValidated:			pPlayer.ProfValidated,
				RulesAccepted:			pPlayer.RulesAccepted,
				Twitch:					pPlayer.Twitch,
				CustomMapsConfirmed:	pPlayer.CustomMapsConfirmed,
				LastCampaignsPlayed:	strings.Join(pPlayer.LastCampaignsPlayed, "|"),
				});
			players.MuPlayers.Unlock();
		}
	}
}
