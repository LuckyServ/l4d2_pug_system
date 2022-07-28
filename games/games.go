package games

import (
	"sync"
	"strconv"
	"time"
	"../players"
)

type EntGame struct {
	ID					string
	CreatedAt			int64 //milliseconds
	PlayersUnpaired		[]*players.EntPlayer
	PlayersA			[]*players.EntPlayer
	PlayersB			[]*players.EntPlayer
	GameConfig			string
	CampaignName		string
	Maps				[]string
	State				int
}

const ( //game states
	StateDummy int = iota
	StateCreating
	StateCreated
	CampaignChosen
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