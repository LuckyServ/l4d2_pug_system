package lobby

import (
	//"fmt"
	"../players"
	"sync"
	"time"
)

type EntLobby struct {
	ID				string
	MmrMin			int
	MmrMax			int
	CreatedAt		int64 //milliseconds
	Players			[]*players.EntPlayer
	GameConfig		string
	ReadyUpState	bool
}
var MapLobbies map[string]*EntLobby = make(map[string]*EntLobby);
var ArrayLobbies []*EntLobby; //duplicate of MapLobbies, for faster iterating
var MuLobbies sync.Mutex;
var I64LastLobbyListUpdate int64;


func Create(iMmr int) (string, int) {

	//calculate mmr limits
	iMmrMin, iMmrMax, errMmrLimits := CalcMmrLimits(iMmr);
	if (errMmrLimits != nil) {
		return "", 1; //error calculating mmr range, shouldn't ever happen
	}

	//get time
	i64CurTime := time.Now().UnixMilli();

	//Choose Confogl config
	sConfoglConfig := ChooseConfoglConfig(iMmr);


	//Write lobby
	MuLobbies.Lock();
	sLobbyID := GenerateID();

	pLobby := &EntLobby{
		ID:					sLobbyID,
		MmrMin:				iMmrMin,
		MmrMax:				iMmrMax,
		CreatedAt:			i64CurTime,
		GameConfig:			sConfoglConfig,
	};
	MapLobbies[sLobbyID] = pLobby;
	ArrayLobbies = append(ArrayLobbies, pLobby);
	I64LastLobbyListUpdate = time.Now().UnixMilli();
	MuLobbies.Unlock();

	return sLobbyID, 0;
}

func Join(pPlayer *players.EntPlayer, sLobbyID string) { //MuPlayers must be locked outside. MuLobbies will be locked inside.
	MuLobbies.Lock();

	pLobby := MapLobbies[sLobbyID];
	pLobby.Players = append(pLobby.Players, pPlayer);

	i64CurTime := time.Now().UnixMilli();
	I64LastLobbyListUpdate = i64CurTime;

	MuLobbies.Unlock();

	pPlayer.IsInLobby = true;
	pPlayer.LastChanged = i64CurTime;
	players.I64LastPlayerlistUpdate = i64CurTime;
}