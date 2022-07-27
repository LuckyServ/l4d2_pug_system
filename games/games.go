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
	Maps				[]string
	State				int
}

const ( //game states
	StateDummy int = iota
	StateCreating
	StatePinging
	StateLookingForServ
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