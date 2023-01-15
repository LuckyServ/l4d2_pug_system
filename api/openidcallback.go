package api

import (
	"github.com/gin-gonic/gin"
	"../settings"
	"io/ioutil"
	"fmt"
	"time"
	"../players"
	"../bans"
	"github.com/yohcop/openid-go"
	"regexp"
	"github.com/antchfx/xmlquery"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"../smurf"
	"crypto/sha256"
)

type NoOpDiscoveryCache struct{};
var nonceStore = openid.NewSimpleNonceStore();
var discoveryCache = &NoOpDiscoveryCache{};
func (n *NoOpDiscoveryCache) Put(id string, info openid.DiscoveredInfo) {}
func (n *NoOpDiscoveryCache) Get(id string) openid.DiscoveredInfo {
	return nil;
}
var mapIPs = make(map[string]int);
var MuAuth sync.RWMutex;

//Limit authorizations per hour per IP
func AuthRatelimits() {
	go TwitchAuthRatelimits();
	for {
		time.Sleep(3600 * time.Second); //1 hour
		MuAuth.Lock();
		mapIPs = make(map[string]int);
		MuAuth.Unlock();
	}
}

func GetAuthCount(sClientIP string) int {
	MuAuth.RLock();
	iCount, bExists := mapIPs[sClientIP];
	MuAuth.RUnlock();
	if (bExists) {
		return iCount;
	} else {
		return 0;
	}
}

func IncreaseAuthCount(sClientIP string) {
	MuAuth.Lock();
	iCount, bExists := mapIPs[sClientIP];
	if (bExists) {
		mapIPs[sClientIP] = iCount + 1;
	} else {
		mapIPs[sClientIP] = 1;
	}
	MuAuth.Unlock();
}


func HttpReqOpenID(c *gin.Context) {

	//Ratelimits
	sClientIP := c.ClientIP();
	iCount := GetAuthCount(sClientIP);
	if (iCount >= settings.AuthPerHour) {
		c.String(200, "Too many authorization requests. Wait an hour before trying again.");
		return;
	}
	IncreaseAuthCount(sClientIP);

	//Get parameters
	mapParameters := c.Request.URL.Query();
	//fmt.Printf("%v\n", mapParameters);
	//Homepage
	sHomepage, errHomepage := c.Cookie("home_page");
	if (errHomepage != nil || sHomepage == "") {
		sHomepage = "https://"+settings.HomeDomain;
	}

	//Check if Steam url valid
	if _, ok := mapParameters["openid.op_endpoint"]; !ok {
		c.Redirect(303, sHomepage);
		return;
	}
	if (len(mapParameters["openid.op_endpoint"]) <= 0) {
		c.Redirect(303, sHomepage);
		return;
	}
	if (mapParameters["openid.op_endpoint"][0] != "https://steamcommunity.com/openid/login") {
		c.Redirect(303, sHomepage);
		return;
	}

	//Validate auth request with Steam
	sReqString := "?dummy=1";
	for sKey, arValues := range mapParameters {
		if (len(arValues) > 0 && sKey != "openid.mode") {
			sReqString = fmt.Sprintf("%s&%s=%s", sReqString, sKey, url.QueryEscape(arValues[0]));
		}
	}
	fullURL := "https://"+settings.BackendDomain + c.Request.URL.Path + sReqString;
	id, err := openid.Verify(fullURL, discoveryCache, nonceStore);
	if (err != nil) {
		c.Redirect(303, sHomepage);
		return;
	}
	vRegEx := regexp.MustCompile(`[0-9]{17}`);
	bySteamID64 := vRegEx.Find([]byte(id));
	if (bySteamID64 == nil) {
		c.Redirect(303, sHomepage);
	}

	//Here is authorized SteamID64
	sSteamID64 := string(bySteamID64);

	//Get nickname
	sNickname := "unknown";
	sAvatarSmall := "https://l4d2center.com/blank_avatar_small.jpg";
	sAvatarBig := "https://l4d2center.com/blank_avatar_big.jpg";
	clientSteam := http.Client{
		Timeout: 15 * time.Second,
	}
	respSteam, errSteam := clientSteam.Get("https://steamcommunity.com/profiles/"+sSteamID64+"/?xml=1");
	if (errSteam != nil) {
		c.Redirect(303, sHomepage);
		return;
	}
	defer respSteam.Body.Close();
	if (respSteam.StatusCode != 200) {
		c.Redirect(303, sHomepage);
		return;
	}
	byResult, errBody := ioutil.ReadAll(respSteam.Body);
	sResult := string(byResult);
	if (errBody != nil || sResult == "") {
		c.Redirect(303, sHomepage);
		return;
	}
	doc, errXML := xmlquery.Parse(strings.NewReader(sResult));
	if (errXML != nil) {
		c.Redirect(303, sHomepage);
		return;
	}
	root := xmlquery.FindOne(doc, "//profile");
	if n := root.SelectElement("//steamID"); n != nil {
		sNickname = n.InnerText();
	}
	if n := root.SelectElement("//avatarMedium"); n != nil {
		sAvatarSmall = n.InnerText();
	}
	if n := root.SelectElement("//avatarFull"); n != nil {
		sAvatarBig = n.InnerText();
	}

	//Add auth to the database
	sSessionID := players.AddPlayerAuth(sSteamID64, base64.StdEncoding.EncodeToString([]byte(sNickname)), sAvatarSmall, sAvatarBig);
	bans.ChanSearchBan <- sSteamID64;

	//smurf
	sCookieUniqueKey, _ := c.Cookie("auth2");
	go smurf.AnnounceIPAndKey(sSteamID64, c.ClientIP(), sNickname, sCookieUniqueKey);

	//Set cookie
	c.SetCookie("session_id", sSessionID, 2592000, "/", "", true, true);

	//Redirect to home page
	by32Buffer := sha256.Sum256([]byte(sSessionID));
	byBuffer := by32Buffer[:];
	sCSRF := base64.StdEncoding.EncodeToString(byBuffer);
	c.Redirect(303, sHomepage + "?auth3=" + url.QueryEscape(sCSRF));
	//fmt.Printf("New auth: %s, %s, %s\n", sSteamID64, sNickname, sCSRF);
}
