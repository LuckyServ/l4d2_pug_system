package api

import (
	"github.com/gin-gonic/gin"
	"../settings"
	"io/ioutil"
	"fmt"
	"time"
	"../players"
	"github.com/yohcop/openid-go"
	"regexp"
	"github.com/antchfx/xmlquery"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
)

type NoOpDiscoveryCache struct{};
var nonceStore = openid.NewSimpleNonceStore();
var discoveryCache = &NoOpDiscoveryCache{};
func (n *NoOpDiscoveryCache) Put(id string, info openid.DiscoveredInfo) {}
func (n *NoOpDiscoveryCache) Get(id string) openid.DiscoveredInfo {
	return nil;
}
var chClear chan bool = make(chan bool);
var chGetCountInList chan string = make(chan string);
var chCountInListResult chan int = make(chan int);
var chAddToList chan string = make(chan string);

//Limit authorizations per hour per IP
func AuthRatelimits() {
	mapIPs := make(map[string]int);
	go func() {
		for {
			time.Sleep(3600 * time.Second); //1 hour
			chClear <- true;
		}
	}();
	for {
		select {
		case <-chClear:
			mapIPs = make(map[string]int);
		case sIP := <-chGetCountInList:
			iCount, bExists := mapIPs[sIP];
			if (bExists) {
				chCountInListResult <- iCount;
			} else {
				chCountInListResult <- 0;
			}
		case sIP := <-chAddToList:
			iCount, bExists := mapIPs[sIP];
			if (bExists) {
				mapIPs[sIP] = iCount + 1;
			} else {
				mapIPs[sIP] = 1;
			}
		}
	}
}


func HttpReqOpenID(c *gin.Context) {
	mapParameters := c.Request.URL.Query();

	//Ratelimits
	sClientIP := c.ClientIP();
	chGetCountInList <- sClientIP;
	iCount := <-chCountInListResult;
	if (iCount >= settings.AuthPerHour) {
		c.String(200, "Too many authorization requests. Wait an hour before trying again.");
		return;
	}
	chAddToList <- sClientIP;

	//Check if Steam url valid
	if _, ok := mapParameters["openid.op_endpoint"]; !ok {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	if (len(mapParameters["openid.op_endpoint"]) <= 0) {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	if (mapParameters["openid.op_endpoint"][0] != "https://steamcommunity.com/openid/login") {
		c.Redirect(303, "https://"+settings.HomeDomain);
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
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	vRegEx := regexp.MustCompile(`[0-9]{17}`);
	bySteamID64 := vRegEx.Find([]byte(id));
	if (bySteamID64 == nil) {
		c.Redirect(303, "https://"+settings.HomeDomain);
	}

	//Here is authorized SteamID64
	sSteamID64 := string(bySteamID64);

	//Get nickname
	sNickname := "unknown";
	clientSteam := http.Client{
		Timeout: 15 * time.Second,
	}
	respSteam, errSteam := clientSteam.Get("https://steamcommunity.com/profiles/"+sSteamID64+"/?xml=1");
	if (errSteam != nil) {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	defer respSteam.Body.Close();
	if (respSteam.StatusCode != 200) {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	byResult, errBody := ioutil.ReadAll(respSteam.Body);
	sResult := string(byResult);
	if (errBody != nil || sResult == "") {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	doc, errXML := xmlquery.Parse(strings.NewReader(sResult));
	if (errXML != nil) {
		c.Redirect(303, "https://"+settings.HomeDomain);
		return;
	}
	root := xmlquery.FindOne(doc, "//profile");
	if n := root.SelectElement("//steamID"); n != nil {
		sNickname = n.InnerText();
	}

	//Add auth to the database
	sSessionID := players.AddPlayerAuth(sSteamID64, base64.StdEncoding.EncodeToString([]byte(sNickname)));

	//Set cookie
	c.SetCookie("session_id", sSessionID, 2592000, "/", "", true, false);

	//Redirect to home page
	c.Header("Access-Control-Allow-Origin", "*");
	c.Redirect(303, "https://"+settings.HomeDomain);
}
