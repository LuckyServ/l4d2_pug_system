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
	ReadyPlayers	int
	ReadyUpSince	int64 //timestamp of initiation of the readyup state //milliseconds
}
var MapLobbies map[string]*EntLobby = make(map[string]*EntLobby);
var ArrayLobbies []*EntLobby; //duplicate of MapLobbies, for faster iterating
var MuLobbies sync.RWMutex;
var I64LastLobbyListUpdate int64;


func JoinAny(pPlayer *players.EntPlayer) bool { //MuPlayers and MuLobbies must be locked outside
	arLobbies := GetJoinableLobbies(pPlayer.Mmr);
	iSize := len(arLobbies);
	if (iSize == 0) {

		if (Create(pPlayer)) {
			return true;
		}
		return false;

	} else {
		//sort
		if (iSize > 1) {
			bSorted := false;
			for !bSorted {
				bSorted = true;
				for i := 1; i < iSize; i++ {
					if (arLobbies[i].CreatedAt < arLobbies[i - 1].CreatedAt) {
						arLobbies[i], arLobbies[i - 1] = arLobbies[i - 1], arLobbies[i]; //switch
						if (bSorted) {
							bSorted = false;
						}
					}
				}
				if (!bSorted) {
					for i := iSize - 2; i >= 0; i-- {
						if (arLobbies[i].CreatedAt > arLobbies[i + 1].CreatedAt) {
							arLobbies[i], arLobbies[i + 1] = arLobbies[i + 1], arLobbies[i]; //switch
						}
					}
				}
			}
		}
		sLobbyID := arLobbies[0].ID;

		if (Join(pPlayer, sLobbyID)) {
			return true;
		}
		return false;
	}
	return false;
}

func Create(pPlayer *players.EntPlayer) bool { //MuPlayers and MuLobbies must be locked outside

	//Repeat some critical checks
	if (pPlayer.IsInLobby) {
		return false;
	}

	//calculate mmr limits
	iMmrMin, iMmrMax, errMmrLimits := CalcMmrLimits(pPlayer);
	if (errMmrLimits != nil) {
		return false; //error calculating mmr range, shouldn't ever happen
	}

	//Write lobby
	i64CurTime := time.Now().UnixMilli();
	sLobbyID := GenerateID();

	pLobby := &EntLobby{
		ID:					sLobbyID,
		MmrMin:				iMmrMin,
		MmrMax:				iMmrMax,
		CreatedAt:			i64CurTime,
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
	if (pLobby.PlayerCount >= 8) { //hardcoded for 4v4 games
		return false;
	}
	if (pPlayer.Mmr < pLobby.MmrMin || pPlayer.Mmr > pLobby.MmrMax) {
		return false;
	}

	pLobby.Players = append(pLobby.Players, pPlayer);
	pLobby.PlayerCount++;

	i64CurTime := time.Now().UnixMilli();
	if (pLobby.PlayerCount == 8) { //hardcoded for 4v4 games
		pLobby.ReadyUpSince = i64CurTime;
	}

	I64LastLobbyListUpdate = i64CurTime;

	pPlayer.IsInLobby = true;
	pPlayer.IsIdle = false;
	pPlayer.LastLobbyActivity = i64CurTime;
	pPlayer.LobbyID = sLobbyID;
	players.I64LastPlayerlistUpdate = i64CurTime;
	return true;
}

func Leave(pPlayer *players.EntPlayer) bool { //MuPlayers and MuLobbies must be locked outside

	//Repeat some critical checks
	if (!pPlayer.IsInLobby) {
		return false;
	}

	//find lobby
	pLobby, bFound := MapLobbies[pPlayer.LobbyID];
	if (!bFound) {
		return false;
	}
	
	//remove player from array
	bRemoved := false;
	for i := 0; i < pLobby.PlayerCount; i++ {
		if (pLobby.Players[i].SteamID64 == pPlayer.SteamID64) {
			pLobby.Players[i] = pLobby.Players[pLobby.PlayerCount - 1];
			pLobby.Players = pLobby.Players[:(pLobby.PlayerCount - 1)];
			bRemoved = true;
			break;
		}
	}
	if (!bRemoved) {
		return false;
	} else {
		pLobby.PlayerCount--;
	}

	i64CurTime := time.Now().UnixMilli();

	if (pLobby.PlayerCount == 0) { //destroy lobby if it's empty
		delete(MapLobbies, pLobby.ID);
		iMaxL := len(ArrayLobbies);
		for i := 0; i < iMaxL; i++ {
			if (ArrayLobbies[i].ID == pLobby.ID) {
				ArrayLobbies[i] = ArrayLobbies[iMaxL - 1];
				ArrayLobbies = ArrayLobbies[:(iMaxL - 1)];
				break;
			}
		}
	} else if (pLobby.PlayerCount == 7) { //Remove ReadyUp traces if the lobby was in ReadyUp state
		if (pPlayer.IsReadyInLobby) {
			pPlayer.IsReadyInLobby = false;
			pLobby.ReadyPlayers--;
		}
		pPlayer.LastFullLobbyLeave = i64CurTime;
		for i := 0; i < pLobby.PlayerCount; i++ {
			if (pLobby.Players[i].IsReadyInLobby) {
				pLobby.Players[i].IsReadyInLobby = false;
				pLobby.ReadyPlayers--;
			}
		}
	}

	I64LastLobbyListUpdate = i64CurTime;

	pPlayer.IsInLobby = false;
	pPlayer.IsIdle = false;
	pPlayer.LastLobbyActivity = i64CurTime;
	pPlayer.LobbyID = "";
	players.I64LastPlayerlistUpdate = i64CurTime;
	return true;
}

func Ready(pPlayer *players.EntPlayer) bool { //MuPlayers and MuLobbies must be locked outside

	//Repeat some critical checks
	if (!pPlayer.IsInLobby) {
		return false;
	} else if (pPlayer.IsReadyInLobby) {
		return false;
	}

	//find lobby
	pLobby, bFound := MapLobbies[pPlayer.LobbyID];
	if (!bFound) {
		return false;
	}

	//Is lobby in readyup state
	if (pLobby.PlayerCount != 8) { //hardcode for 4v4
		return false;
	}

	i64CurTime := time.Now().UnixMilli();

	//Set ready
	pPlayer.IsReadyInLobby = true;
	pLobby.ReadyPlayers++;

	I64LastLobbyListUpdate = i64CurTime;

	return true;
}
