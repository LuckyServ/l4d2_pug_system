package games

import (
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

var MapGames map[string]*EntGame = make(map[string]*EntGame);
var ArrayGames []*EntGame; //duplicate of MapGames, for faster iterating
var MuGames sync.Mutex;

var ChanNewGameID chan string = make(chan string);

func ChannelWatchers() {
	go HandleUniqID();
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
	for _, pPlayer := range pGame.PlayersUnpaired {
		pPlayer.IsInGame = true;
		pPlayer.GameID = pGame.ID;
	}
	players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
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