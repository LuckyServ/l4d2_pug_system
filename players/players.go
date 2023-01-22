package players

import (
	"sync"
	"time"
	"../settings"
	"../utils"
	"../database"
	"./auth"
	"../smurf"
	"strconv"
	"encoding/base64"
	"strings"
)

type EntPlayer struct {
	SteamID64				string
	NicknameBase64			string
	AvatarSmall				string
	AvatarBig				string
	Mmr						int
	MmrUncertainty			float32
	LastGameResult			int //0 - unknown, 1 - draw, 2 - lost, 3 - won
	Access					int //-3 - banned + cant protest, -2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	BannedAt				int64 //unix timestamp in milliseconds
	BanReason				string
	BanAcceptedAt			int64 //unix timestamp in milliseconds
	BanLength				int64 //unix timestamp in milliseconds
	ProfValidated			bool //Steam profile validated
	RulesAccepted			bool //Rules accepted
	LastActivity			int64 //unix timestamp in milliseconds
	IsOnline				bool
	OnlineSince				int64 //unix timestamp in milliseconds
	IsInGame				bool
	IsInQueue				bool
	InQueueSince			int64 //unix timestamp in milliseconds
	LastQueueChanged		int64 //unix timestamp in milliseconds
	IsReadyUpRequested		bool
	IsReadyConfirmed		bool
	NextQueueingAllowed		int64 //unix timestamp in milliseconds
	GameID					string
	LastGameChanged			int64  //unix timestamp in milliseconds
	LastExternalRequest		int64 //Last external api request //unix timestamp in milliseconds
	LastLimitedRequest		int64 //unix timestamp in milliseconds
	LastGameActivity		int64 //Last game activity //unix timestamp in milliseconds
	LastChatMessage			int64 //Last chat message //unix timestamp in milliseconds
	LastTicketActivity		int64 //Last ticket activity //unix timestamp in milliseconds
	GameServerPings			map[string]int
	LastCampaignsPlayed		[]string
	Twitch					string
	DuoWith					string
	DuoOffer				string
	CustomMapsConfirmed		int64 //unix timestamp in milliseconds
}

var MapPlayers map[string]*EntPlayer = make(map[string]*EntPlayer);
var ArrayPlayers []*EntPlayer; //duplicate of MapPlayers, for faster iterating
var MuPlayers sync.RWMutex;

var I64LastPlayerlistUpdate int64;
var iMaxStableMmr, iMinStableMmr, iMmrDiff int;


func Watchers() {
	go WatchOnline();
	go WatchSortPlayers();
}

func UpdatePlayerActivity(sSteamID64 string, sCookieUniqueKey string, sIP string) { //Maps must be locked outside!!!
	if _, ok := MapPlayers[sSteamID64]; !ok {
		return;
	}
	i64CurTime := time.Now().UnixMilli();
	pPlayer := MapPlayers[sSteamID64];
	pPlayer.LastActivity = i64CurTime;
	if (!pPlayer.IsOnline) {
		pPlayer.IsOnline = true;
		pPlayer.OnlineSince = i64CurTime;
		I64LastPlayerlistUpdate = i64CurTime;

		byNickname, _ := base64.StdEncoding.DecodeString(pPlayer.NicknameBase64);
		go smurf.AnnounceIPAndKey(pPlayer.SteamID64, sIP, string(byNickname), sCookieUniqueKey);
	}
}

func RestorePlayers() bool { //no need to lock maps
	arDatabasePlayers := database.RestorePlayers();
	/*if (len(arDatabasePlayers) == 0) {
		return false;
	}*/
	//var iBottomMmr int = 1;
	for _, oDBPlayer := range arDatabasePlayers {
		iAccess := oDBPlayer.Access;
		if (iAccess < 0) {
			iAccess = 0;
		}
		pPlayer := &EntPlayer{
			SteamID64:				oDBPlayer.SteamID64,
			NicknameBase64:			oDBPlayer.NicknameBase64,
			AvatarSmall:			oDBPlayer.AvatarSmall,
			AvatarBig:				oDBPlayer.AvatarBig,
			Mmr:					oDBPlayer.Mmr,
			MmrUncertainty:			oDBPlayer.MmrUncertainty,
			LastGameResult:			oDBPlayer.LastGameResult,
			Access:					iAccess,
			ProfValidated:			oDBPlayer.ProfValidated,
			RulesAccepted:			oDBPlayer.RulesAccepted,
			Twitch:					oDBPlayer.Twitch,
			CustomMapsConfirmed:	oDBPlayer.CustomMapsConfirmed,
			LastCampaignsPlayed:	func()([]string){
				if (oDBPlayer.LastCampaignsPlayed == "") {
					return []string{};
				} else {
					return strings.Split(oDBPlayer.LastCampaignsPlayed, "|");
				}
			}(),
		};
		MapPlayers[oDBPlayer.SteamID64] = pPlayer;
		ArrayPlayers = append(ArrayPlayers, pPlayer);
		/*if (oDBPlayer.ProfValidated && oDBPlayer.Mmr < iBottomMmr && GetMmrGrade(pPlayer) > 0) {
			iBottomMmr = oDBPlayer.Mmr;
		}*/
	}
	/*if (iBottomMmr < 1) {
		iShiftMmr := 1 - iBottomMmr;
		database.ShiftMmr(iShiftMmr);
		return false;
	}*/
	I64LastPlayerlistUpdate = time.Now().UnixMilli();
	return true;
}


func AddPlayerAuth(sSteamID64 string, sNicknameBase64 string, sAvatarSmall string, sAvatarBig string) string {

	//Register player if does not exist
	MuPlayers.Lock();
	if _, ok := MapPlayers[sSteamID64]; !ok {

		pPlayer := &EntPlayer{
			SteamID64:			sSteamID64,
			MmrUncertainty:		settings.DefaultMmrUncertainty,
		};
		MapPlayers[sSteamID64] = pPlayer;
		ArrayPlayers = append(ArrayPlayers, pPlayer);
		I64LastPlayerlistUpdate = time.Now().UnixMilli();

		oDatabasePlayer := database.DatabasePlayer{
			SteamID64:			sSteamID64,
			MmrUncertainty:		settings.DefaultMmrUncertainty,
		};
		MuPlayers.Unlock();

		database.AddPlayer(oDatabasePlayer);

	} else {
		MuPlayers.Unlock();
	}

	sSessionKey, _ := utils.GenerateRandomString(32, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz");
	sSessionKey = sSessionKey+strconv.FormatInt(time.Now().UnixNano(), 10);

	auth.MuSessions.Lock();

	MuPlayers.Lock();
	bEdited := false;
	pPlayer := MapPlayers[sSteamID64];
	if (pPlayer.NicknameBase64 != sNicknameBase64) {
		pPlayer.NicknameBase64 = sNicknameBase64;
		bEdited = true;
	}
	if (pPlayer.AvatarSmall != sAvatarSmall) {
		pPlayer.AvatarSmall = sAvatarSmall;
		bEdited = true;
	}
	if (pPlayer.AvatarBig != sAvatarBig) {
		pPlayer.AvatarBig = sAvatarBig;
		bEdited = true;
	}
	if (bEdited) {
		I64LastPlayerlistUpdate = time.Now().UnixMilli();
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
	}
	MuPlayers.Unlock();

	oSession := auth.EntSession{
		SteamID64:	sSteamID64,
		Since:		time.Now().UnixMilli(),
	};

	auth.MapSessions[sSessionKey] = oSession;
	go database.AddSession(database.DatabaseSession{
		SessionID:			sSessionKey,
		SteamID64:			oSession.SteamID64,
		Since:				oSession.Since,
		});
	auth.MuSessions.Unlock();

	return sSessionKey;
}

func WatchSortPlayers() {
	for {
		time.Sleep(5 * time.Second);
		
		MuPlayers.Lock();
		SortPlayers();
		MuPlayers.Unlock();
	}
}

func SortPlayers() { //Players must be locked outside
	bEdited := false;
	iSize := len(ArrayPlayers);
	if (iSize > 1) {
		bSorted := false;
		for !bSorted {
			bSorted = true;
			for i := 1; i < iSize; i++ {
				if (ArrayPlayers[i].Mmr > ArrayPlayers[i - 1].Mmr) {
					ArrayPlayers[i], ArrayPlayers[i - 1] = ArrayPlayers[i - 1], ArrayPlayers[i]; //switch
					if (!bEdited) {
						bEdited = true;
					}
					if (bSorted) {
						bSorted = false;
					}
				}
			}
			if (!bSorted) {
				for i := iSize - 2; i >= 0; i-- {
					if (ArrayPlayers[i].Mmr < ArrayPlayers[i + 1].Mmr) {
						ArrayPlayers[i], ArrayPlayers[i + 1] = ArrayPlayers[i + 1], ArrayPlayers[i]; //switch
						if (!bEdited) {
							bEdited = true;
						}
					}
				}
			}
		}
	}
	for i := iSize - 1; i >= 0; i-- {
		if (ArrayPlayers[i].MmrUncertainty <= settings.MmrStable) {
			iMinStableMmr = ArrayPlayers[i].Mmr;
			break;
		}
	}
	for i := 0; i < iSize; i++ {
		if (ArrayPlayers[i].MmrUncertainty <= settings.MmrStable) {
			iMaxStableMmr = ArrayPlayers[i].Mmr;
			break;
		}
	}
	if (iMaxStableMmr > iMinStableMmr) {
		iMmrDiff = iMaxStableMmr - iMinStableMmr;
	}
	if (bEdited) {
		I64LastPlayerlistUpdate = time.Now().UnixMilli();
	}
}

func GetMmrGrade(pPlayer *EntPlayer) int { //Players must be locked outside
	if (pPlayer.MmrUncertainty > settings.MmrStable) {
		return 0;
	}
	if (pPlayer.Mmr < iMinStableMmr || pPlayer.Mmr > iMaxStableMmr) {
		return 0;
	}
	if (iMmrDiff < 1) {
		return 0;
	}
	iGrade := (pPlayer.Mmr - iMinStableMmr) / (iMmrDiff / settings.MmrGrades);
	if (iGrade < settings.MmrGrades) {
		iGrade = iGrade + 1;
	}
	return iGrade;
}

func CustomMapsConfirmState(pPlayer *EntPlayer) int { //Players must be locked outside
	if (pPlayer.CustomMapsConfirmed == 0) {
		return 1;
	} else if (settings.NewestCustomMap < pPlayer.CustomMapsConfirmed) {
		return 3;
	}
	return 2;
}