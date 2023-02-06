package streams

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"io/ioutil"
	"github.com/buger/jsonparser"
	"sync"
	"../players"
	"../settings"
	"../utils"
	//"errors"
)

type TwitchStream struct {
	UserLogin			string		`json:"user_login"`
	UserName			string		`json:"user_name"`
	Title				string		`json:"stream_title"`
	Language			string		`json:"stream_language"`
	Thumbnail			string		`json:"stream_thumbnail"`
	ViewersCount		int			`json:"viewers_count"`
}

var ArrayStreams []TwitchStream;
var MuStreams sync.RWMutex;
var I64LastStreamersUpdate int64;

var sTwitchAuthToken string;

func Watchers() {
	go WatchOnline();
}

func WatchOnline() {
	for {
		time.Sleep(30 * time.Second);
		UpdateOnlineStreams();
	}
}

func UpdateOnlineStreams() {
	//https://dev.twitch.tv/docs/api/reference/#get-streams

	var arPrioritizedTwitchIDs, arOtherTwitchIDs []string;
	players.MuPlayers.RLock();
	for _, pPlayer := range players.ArrayPlayers {
		if ((pPlayer.IsInQueue || pPlayer.IsInGame) && pPlayer.Twitch != "") {
			arPrioritizedTwitchIDs = append(arPrioritizedTwitchIDs, pPlayer.Twitch);
		} else if (!pPlayer.IsInQueue && !pPlayer.IsInGame && pPlayer.Access >= 0 && players.GetMmrGrade(pPlayer) > 0 && pPlayer.Twitch != "") {
			arOtherTwitchIDs = append(arOtherTwitchIDs, pPlayer.Twitch);
		}
	}
	players.MuPlayers.RUnlock();

	arStreams := GetTwitchStreams(arPrioritizedTwitchIDs);
	arStreams2 := GetTwitchStreams(arOtherTwitchIDs);
	arStreams = SortStreams(arStreams);
	arStreams2 = SortStreams(arStreams2);
	arStreams = append(arStreams, arStreams2...);

	MuStreams.Lock();
	if (!StreamlistEqual(arStreams)) {
		ArrayStreams = arStreams;
		I64LastStreamersUpdate = time.Now().UnixMilli();
	}
	MuStreams.Unlock();
}

func GetTwitchStreams(arTwitchIDs []string) ([]TwitchStream) {
	if (sTwitchAuthToken == "") {
		if (RequestTwitchAuthToken()) {
			arStreams, iRespCode := RequestTwitchStreams(arTwitchIDs);
			if (iRespCode == 200) {
				return arStreams;
			}
		}
	} else {
		arStreams, iRespCode := RequestTwitchStreams(arTwitchIDs);
		if (iRespCode == 200) {
			return arStreams;
		} else if (iRespCode == 401) {
			if (RequestTwitchAuthToken()) {
				arStreams2, iRespCode2 := RequestTwitchStreams(arTwitchIDs);
				if (iRespCode2 == 200) {
					return arStreams2;
				}
			}
		}
	}
	return []TwitchStream{};
}

func RequestTwitchAuthToken() bool {
	clientOAuth := http.Client{
		Timeout: 10 * time.Second,
	}
	data := url.Values{};
	data.Set("client_id", settings.TwitchClientID);
	data.Set("client_secret", settings.TwitchSecret);
	data.Set("grant_type", "client_credentials");
	encodedData := data.Encode();
	reqOAuth, _ := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", strings.NewReader(encodedData));
	reqOAuth.Header.Set("Content-Type", "application/x-www-form-urlencoded");
	respOAuth, errOAuth := clientOAuth.Do(reqOAuth);
	if (errOAuth != nil) {
		return false;
	}
	defer respOAuth.Body.Close();
	if (respOAuth.StatusCode != 200) {
		return false;
	}
	byResBody, errResBody := ioutil.ReadAll(respOAuth.Body);
	if (errResBody != nil) {
		return false;
	}
	sAccesToken, errAccesToken := jsonparser.GetString(byResBody, "access_token");
	if (errAccesToken != nil || sAccesToken == "") {
		return false;
	} else {
		sTwitchAuthToken = sAccesToken;
		return true;
	}
	return false;
}

func RequestTwitchStreams(arTwitchIDs []string) ([]TwitchStream, int) {
	var arStreams []TwitchStream;

	for {
		iRemainedIDs := len(arTwitchIDs);
		if (iRemainedIDs == 0) {
			break;
		}
		iToCheck := utils.MinValInt(iRemainedIDs, 100);




		clientStreams := http.Client{
			Timeout: 10 * time.Second,
		}
		var sUrl string = "https://api.twitch.tv/helix/streams?type=live&game_id=24193&";
		for i := 0; i < iToCheck; i++ {
			sUrl = fmt.Sprintf("%suser_id=%s&", sUrl, arTwitchIDs[i]);
		}
		arTwitchIDs = arTwitchIDs[iToCheck:];
		reqStreams, _ := http.NewRequest("GET", sUrl, nil);
		reqStreams.Header.Set("Authorization", "Bearer "+sTwitchAuthToken);
		reqStreams.Header.Set("Client-Id", settings.TwitchClientID);
		respStreams, errStreams := clientStreams.Do(reqStreams);
		if (errStreams != nil) {
			return arStreams, 0;
		} else {
			if (respStreams.StatusCode == 200) {
				byResBody, errResBody := ioutil.ReadAll(respStreams.Body);
				if (errResBody != nil) {
					respStreams.Body.Close();
					return arStreams, 0;
				}

				jsonparser.ArrayEach(byResBody, func(valueStream []byte, dataType jsonparser.ValueType, offset int, err error) {
					sUserLogin, _ := jsonparser.GetString(valueStream, "user_login");
					sUserName, _ := jsonparser.GetString(valueStream, "user_name");
					sTitle, _ := jsonparser.GetString(valueStream, "title");
					sLanguage, _ := jsonparser.GetString(valueStream, "language");
					sThumbnail, _ := jsonparser.GetString(valueStream, "thumbnail_url");
					i64ViewersCount, _ := jsonparser.GetInt(valueStream, "viewer_count");
					arStreams = append(arStreams, TwitchStream{
						UserLogin:			sUserLogin,
						UserName:			sUserName,
						Title:				sTitle,
						Language:			sLanguage,
						Thumbnail:			sThumbnail,
						ViewersCount:		int(i64ViewersCount),
					});
				}, "data");
			

			} else {
				iStatusCode := respStreams.StatusCode;
				respStreams.Body.Close();
				return arStreams, iStatusCode;
			}
			respStreams.Body.Close();
		}
	}

	return arStreams, 200;
}

func StreamlistEqual(arStreams []TwitchStream) bool {
	iLenStreams := len(ArrayStreams);
	if (iLenStreams != len(arStreams)) {
		return false;
	}
	for i := 0; i < iLenStreams; i++ {
		if (ArrayStreams[i] != arStreams[i]) {
			return false
		}
	}
	return true;
}

func SortStreams(arStreams []TwitchStream) []TwitchStream {
	iSize := len(arStreams);
	if (iSize > 1) {
		bSorted := false;
		for !bSorted {
			bSorted = true;
			for i := 1; i < iSize; i++ {
				if (arStreams[i].ViewersCount > arStreams[i - 1].ViewersCount) {
					arStreams[i], arStreams[i - 1] = arStreams[i - 1], arStreams[i]; //switch
					if (bSorted) {
						bSorted = false;
					}
				}
			}
			if (!bSorted) {
				for i := iSize - 2; i >= 0; i-- {
					if (arStreams[i].ViewersCount < arStreams[i + 1].ViewersCount) {
						arStreams[i], arStreams[i + 1] = arStreams[i + 1], arStreams[i]; //switch
					}
				}
			}
		}
	}
	return arStreams;
}