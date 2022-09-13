package players

import (
	"sync"
	"time"
	"../settings"
	"../utils"
	"../database"
	"./auth"
	"strconv"
)

type EntPlayer struct {
	SteamID64			string
	NicknameBase64		string
	Mmr					int
	MmrUncertainty		float32
	Access				int //-2 - completely banned, -1 - chat banned, 0 - regular player, 1 - behaviour moderator, 2 - cheat moderator, 3 - behaviour+cheat moderator, 4 - full admin access
	BanReason			string
	BanAcceptedAt		int64 //unix timestamp in milliseconds
	BanLength			int64 //unix timestamp in milliseconds
	ProfValidated		bool //Steam profile validated
	RulesAccepted		bool //Rules accepted
	LastActivity		int64 //unix timestamp in milliseconds
	IsOnline			bool
	IsIdle				bool
	OnlineSince			int64 //unix timestamp in milliseconds
	IsInGame			bool
	IsInLobby			bool
	IsAutoSearching		bool
	AutoSearchingSince	int64 //unix timestamp in milliseconds
	LobbyID				string
	GameID				string
	LastGameChanged		int64  //unix timestamp in milliseconds
	IsReadyInLobby		bool
	LastSteamRequest	int64 //Last steam api request //unix timestamp in milliseconds
	LastLobbyActivity	int64 //Last lobby activity //unix timestamp in milliseconds
	LastGameActivity	int64 //Last game activity //unix timestamp in milliseconds
	LastChatMessage		int64 //Last chat message //unix timestamp in milliseconds
	GameServerPings		map[string]int
}

var MapPlayers map[string]*EntPlayer = make(map[string]*EntPlayer);
var ArrayPlayers []*EntPlayer; //duplicate of MapPlayers, for faster iterating
var MuPlayers sync.RWMutex;

var I64LastPlayerlistUpdate int64;

func Watchers() {
	go WatchOnline();
	go SortPlayers();
}

func UpdatePlayerActivity(sSteamID64 string) { //Maps must be locked outside!!!
	if _, ok := MapPlayers[sSteamID64]; !ok {
		return;
	}
	i64CurTime := time.Now().UnixMilli();
	pPlayer := MapPlayers[sSteamID64];
	pPlayer.LastActivity = i64CurTime;
	if (!pPlayer.IsOnline) {
		pPlayer.IsOnline = true;
		pPlayer.OnlineSince = i64CurTime;
		pPlayer.IsIdle = false;
		I64LastPlayerlistUpdate = i64CurTime;
	}
}

func RestorePlayers() bool { //no need to lock maps
	arDatabasePlayers := database.RestorePlayers();
	/*if (len(arDatabasePlayers) == 0) {
		return false;
	}*/
	for _, oDBPlayer := range arDatabasePlayers {
		pPlayer := &EntPlayer{
			SteamID64:			oDBPlayer.SteamID64,
			NicknameBase64:		oDBPlayer.NicknameBase64,
			Mmr:				oDBPlayer.Mmr,
			MmrUncertainty:		oDBPlayer.MmrUncertainty,
			Access:				oDBPlayer.Access,
			ProfValidated:		oDBPlayer.ProfValidated,
			RulesAccepted:		oDBPlayer.RulesAccepted,
		};
		MapPlayers[oDBPlayer.SteamID64] = pPlayer;
		ArrayPlayers = append(ArrayPlayers, pPlayer);
	}
	I64LastPlayerlistUpdate = time.Now().UnixMilli();
	return true;
}


func AddPlayerAuth(sSteamID64 string, sNicknameBase64 string) string {

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
	pPlayer := MapPlayers[sSteamID64];
	if (pPlayer.NicknameBase64 != sNicknameBase64) {

		pPlayer.NicknameBase64 = sNicknameBase64;
		I64LastPlayerlistUpdate = time.Now().UnixMilli();

		go database.UpdatePlayer(database.DatabasePlayer{
			SteamID64:			pPlayer.SteamID64,
			NicknameBase64:		pPlayer.NicknameBase64,
			Mmr:				pPlayer.Mmr,
			MmrUncertainty:		pPlayer.MmrUncertainty,
			Access:				pPlayer.Access,
			ProfValidated:		pPlayer.ProfValidated,
			RulesAccepted:		pPlayer.RulesAccepted,
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

func SortPlayers() {
	for {
		time.Sleep(5 * time.Second);

		bEdited := false;
		MuPlayers.Lock();
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
		MuPlayers.Unlock();
		if (bEdited) {
			I64LastPlayerlistUpdate = time.Now().UnixMilli();
		}
	}
}