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
	PlayerCount		int
	GameConfig		string
	ReadyUpState	bool
}
var MapLobbies map[string]*EntLobby = make(map[string]*EntLobby);
var ArrayLobbies []*EntLobby; //duplicate of MapLobbies, for faster iterating
var MuLobbies sync.Mutex;
var I64LastLobbyListUpdate int64;


func Create(pPlayer *players.EntPlayer) (bool) { //MuPlayers and MuLobbies must be locked outside

	//Repeat some critical checks
	if (pPlayer.IsInLobby) {
		return false;
	}

	//calculate mmr limits
	iMmrMin, iMmrMax, errMmrLimits := CalcMmrLimits(pPlayer);
	if (errMmrLimits != nil) {
		return false; //error calculating mmr range, shouldn't ever happen
	}

	//Choose Confogl config
	sConfoglConfig := ChooseConfoglConfig(pPlayer.Mmr);


	//Write lobby
	i64CurTime := time.Now().UnixMilli();
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
	pLobby.Players = append(pLobby.Players, pPlayer); //join lobby
	pLobby.PlayerCount++;
	I64LastLobbyListUpdate = i64CurTime;
	pPlayer.IsInLobby = true;
	pPlayer.IsIdle = false;
	pPlayer.LastLobbyActivity = i64CurTime;
	pPlayer.LobbyID = sLobbyID;
	pPlayer.LastChanged = i64CurTime;
	players.I64LastPlayerlistUpdate = i64CurTime;

	return true;
}

func Join(pPlayer *players.EntPlayer, sLobbyID string) bool { //MuPlayers and MuLobbies must be locked outside

	pLobby, bExists := MapLobbies[sLobbyID];

	//Repeat some critical checks
	if (!bExists) {
		return false;
	}
	if (pPlayer.IsInLobby) {
		return false;
	}
	if (len(pLobby.Players) >= 8) { //hardcoded for 4v4 games
		return false;
	}
	if (pPlayer.Mmr < pLobby.MmrMin || pPlayer.Mmr > pLobby.MmrMax) {
		return false;
	}

	pLobby.Players = append(pLobby.Players, pPlayer);
	pLobby.PlayerCount++;

	i64CurTime := time.Now().UnixMilli();
	I64LastLobbyListUpdate = i64CurTime;

	pPlayer.IsInLobby = true;
	pPlayer.IsIdle = false;
	pPlayer.LastLobbyActivity = i64CurTime;
	pPlayer.LobbyID = sLobbyID;
	pPlayer.LastChanged = i64CurTime;
	players.I64LastPlayerlistUpdate = i64CurTime;
	return true;
}

func Leave(pPlayer *players.EntPlayer) bool { //MuPlayers and MuLobbies must be locked outside

	//Repeat some critical checks
	if (!pPlayer.IsInLobby) {
		return false;
	}

	//find lobby
	var pLobby *EntLobby;
	bFound := false;
	for _, pLobbyBuffer := range ArrayLobbies {
		for _, pPlayerBuffer := range pLobbyBuffer.Players {
			if (pPlayerBuffer.SteamID64 == pPlayer.SteamID64) {
				bFound = true;
				break;
			}
		}
		if (bFound) {
			pLobby = pLobbyBuffer;
			break;
		}
	}
	if (!bFound) {
		return false;
	}
	
	//remove player from array
	iMaxLP := len(pLobby.Players);
	bRemoved := false;
	for i := 0; i < iMaxLP; i++ {
		if (pLobby.Players[i].SteamID64 == pPlayer.SteamID64) {
			pLobby.Players[i] = pLobby.Players[iMaxLP - 1];
			pLobby.Players = pLobby.Players[:(iMaxLP - 1)];
			bRemoved = true;
			break;
		}
	}
	if (!bRemoved) {
		return false;
	} else {
		pLobby.PlayerCount--;
	}

	//destroy lobby if it's empty
	if (len(pLobby.Players) == 0) {
		delete(MapLobbies, pLobby.ID);
		iMaxL := len(ArrayLobbies);
		for i := 0; i < iMaxL; i++ {
			if (ArrayLobbies[i].ID == pLobby.ID) {
				ArrayLobbies[i] = ArrayLobbies[iMaxL - 1];
				ArrayLobbies = ArrayLobbies[:(iMaxL - 1)];
				break;
			}
		}
	}

	i64CurTime := time.Now().UnixMilli();
	I64LastLobbyListUpdate = i64CurTime;

	pPlayer.IsInLobby = false;
	pPlayer.IsIdle = false;
	pPlayer.LastLobbyActivity = i64CurTime;
	pPlayer.LobbyID = "";
	pPlayer.LastChanged = i64CurTime;
	players.I64LastPlayerlistUpdate = i64CurTime;
	return true;
}
