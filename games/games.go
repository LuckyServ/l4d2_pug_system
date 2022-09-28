package games

import (
	"fmt"
	"sync"
	"strconv"
	"time"
	"../players"
	"../settings"
	"../rating"
)

type EntGame struct {
	ID					string
	CreatedAt			int64 //milliseconds
	PlayersUnpaired		[]*players.EntPlayer
	PlayersA			[]*players.EntPlayer
	PlayersB			[]*players.EntPlayer
	GameConfig			settings.ConfoglConfig
	CampaignName		string
	Maps				[]string
	State				int
	ServerIP			string
	MmrMin				int
	MmrMax				int
	GameResult			rating.EntGameResult
	ReceiverFullRUP		chan bool
	ReceiverReadyList	chan []string
	ReceiverResult		chan rating.EntGameResult
}

const ( //game states
	StateDummy int = iota
	StateCreating
	StateCreated
	StateCampaignChosen
	StateTeamsPicked
	StateWaitPings
	StateSelectServer
	StateNoServers
	StateWaitPlayersJoin
	StateReadyUpExpired
	StateGameProceeds
	StateGameEnded
)

var MapGameStatus map[int]string;

var MapGames map[string]*EntGame = make(map[string]*EntGame);
var ArrayGames []*EntGame; //duplicate of MapGames, for faster iterating
var MuGames sync.RWMutex;

var ChanNewGameID chan string = make(chan string);

func Watchers() {
	go HandleUniqID();
	go CheckVersion();

	//init game statuses text for response
	MapGameStatus = map[int]string{
		StateCreating:			"Creating game",
		StateCreated:			"Game created",
		StateCampaignChosen:	"Campaign selected",
		StateTeamsPicked:		"Players paired",
		StateWaitPings:			"Pinging servers",
		StateSelectServer:		"Selecting best available server",
		StateNoServers:			"No free servers available. If no server found in 5 minutes, the game ends.",
		StateWaitPlayersJoin:	"The server is ready. You have "+fmt.Sprintf("%d", settings.FirstReadyUpExpire / 60)+" minutes to join the server and Ready Up.",
		StateReadyUpExpired:	"Some players failed to Ready Up in time",
		StateGameProceeds:		"In game",
		StateGameEnded:			"Game ended",
	}
}

func HandleUniqID() {
	for {
		select {
			case ChanNewGameID <- func()(string) {
				return strconv.FormatInt(time.Now().UnixNano(), 10);
			}():
		}
		time.Sleep(1 * time.Nanosecond);
	}
}


func Create(pGame *EntGame) { //MuGames and MuPlayers must be locked outside
	MapGames[pGame.ID] = pGame;
	ArrayGames = append(ArrayGames, pGame);
	i64CurTime := time.Now().UnixMilli();
	var iMmrSum int;
	for _, pPlayer := range pGame.PlayersUnpaired {
		pPlayer.IsInGame = true;
		pPlayer.GameID = pGame.ID;
		pPlayer.IsIdle = false;
		pPlayer.LastGameActivity = i64CurTime;
		iMmrSum = iMmrSum + pPlayer.Mmr;
	}
	pGame.GameConfig = ChooseConfoglConfig(iMmrSum / 8);
	players.I64LastPlayerlistUpdate = i64CurTime;
}

func Destroy(pGame *EntGame) { //MuGames and MuPlayers must be locked outside
	delete(MapGames, pGame.ID);
	iMaxG := len(ArrayGames);
	for i := 0; i < iMaxG; i++ {
		if (ArrayGames[i].ID == pGame.ID) {
			ArrayGames[i] = ArrayGames[iMaxG - 1];
			ArrayGames = ArrayGames[:(iMaxG - 1)];
			break;
		}
	}
	i64CurTime := time.Now().UnixMilli();
	for _, pPlayer := range pGame.PlayersUnpaired {
		pPlayer.IsInGame = false;
		pPlayer.GameID = "";
		pPlayer.LastGameChanged = i64CurTime;
		pPlayer.IsIdle = false;
		pPlayer.LastGameActivity = i64CurTime;
	}
	players.I64LastPlayerlistUpdate = i64CurTime;
}

func SetLastUpdated(arPlayers []*players.EntPlayer) { //Players must be locked outside
	i64CurTime := time.Now().UnixMilli();
	for _, pPlayer := range arPlayers {
		pPlayer.LastGameChanged = i64CurTime;
	}
}

func GetGameByIP(sIP string) (*EntGame) { //Games must be locked outside
	if (sIP == "") {
		return nil;
	}
	for _, pGame := range ArrayGames {
		if (pGame.ServerIP == sIP) {
			return pGame;
		}
	}
	return nil;
}

func ChooseConfoglConfig(iMmr int) (settings.ConfoglConfig) {
	if (settings.BrokenMode) {
		return settings.ConfoglConfig{
			CodeName:		"default",
			Name:			"Default",
			MmrMax:			2000000000,
		};
	}
	iLen := len(settings.ArrayConfoglConfigsMmrs);
	if (iLen == 1) {
		return settings.MapConfoglConfigs[settings.ArrayConfoglConfigsMmrs[0]];
	}
	for i := 0; i < iLen; i++ {
		if (iMmr < settings.ArrayConfoglConfigsMmrs[i]) {
			return settings.MapConfoglConfigs[settings.ArrayConfoglConfigsMmrs[i]];
		}
	}
	return settings.ConfoglConfig{
		CodeName:		"zonemod",
		Name:			"nani?",
		MmrMax:			2000000000,
	}; //shouldn't happen
}
