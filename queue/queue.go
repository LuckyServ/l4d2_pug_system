package queue

import (
	"../players"
	"time"
	"errors"
)

var arQueue []*players.EntPlayer;
var NewGamesBlocked bool;
var IPlayersCount int;
var BIsInReadyUp bool;
var i64InReadyUpSince int64;
var PLongestWaitPlayer *players.EntPlayer;
var pPlayerReadyUpReason *players.EntPlayer;
var IReadyPlayers int;
var bWaitingForSinglePlayer bool;

var MapDuoOffers = make(map[string]*players.EntPlayer, 0);


func OfferDuo(pPlayer *players.EntPlayer) { //Players must be locked outside
	if (pPlayer.DuoOffer != "") {
		delete(MapDuoOffers, pPlayer.DuoOffer);
	}
	pPlayer.DuoOffer = <-GenInviteCode;
	MapDuoOffers[pPlayer.DuoOffer] = pPlayer;
	players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
}

func AcceptDuo(pPlayer *players.EntPlayer, sInviteCode string) error { //Players must be locked outside

	pPlayer2, bInviteFound := MapDuoOffers[sInviteCode];
	if (bInviteFound && pPlayer2 != nil) {
		if (pPlayer2.IsInQueue || pPlayer2.IsInGame || pPlayer2.DuoOffer == "") {
			return errors.New("Very bad error happened, go slap the admin");
		}
		if (pPlayer2 == pPlayer) {
			return errors.New("You cannot accept your own Invite Code");
		}
		if (!pPlayer2.IsOnline) {
			return errors.New("Your friend is not online");
		}
		if (pPlayer2.Access <= -2) {
			return errors.New("Your friend is banned");
		}
		if (players.GetMmrGrade(pPlayer2) >= 6 && players.CustomMapsConfirmState(pPlayer2) != 3) {
			return errors.New("Your friend needs to download/install custom maps before playing");
		}

		pPlayer.DuoWith = pPlayer2.SteamID64;
		pPlayer2.DuoWith = pPlayer.SteamID64;

		Join(pPlayer);
		Join(pPlayer2);

		return nil;
	}
	return errors.New("Invite Code not valid");
}

func CancelDuo(pPlayer *players.EntPlayer) { //Players must be locked outside
	delete(MapDuoOffers, pPlayer.DuoOffer);
	pPlayer.DuoOffer = "";
	players.I64LastPlayerlistUpdate = time.Now().UnixMilli();
}


func Join(pPlayer *players.EntPlayer) { //Players must be locked outside

	i64CurTime := time.Now().UnixMilli();
	arQueue = append(arQueue, pPlayer);
	IPlayersCount++;
	pPlayer.IsInQueue = true;

	var iOnline int;
	for _, pPlayer := range players.ArrayPlayers {
		if ((pPlayer.IsOnline || pPlayer.IsInGame || pPlayer.IsInQueue) && pPlayer.ProfValidated && pPlayer.RulesAccepted && pPlayer.Access >= -1/*not banned*/) {
			iOnline++;
		}
	}

	if (iOnline >= 60) {
		pPlayer.InQueueSince = i64CurTime;
	} else {
		pPlayer.InQueueSince = i64CurTime - I64MaxQueueWait;
	}
	pPlayer.IsReadyUpRequested = false;
	pPlayer.IsReadyConfirmed = false;
	if (IPlayersCount == 1) {
		PLongestWaitPlayer = pPlayer;
	}
	if (bWaitingForSinglePlayer && pPlayer.DuoWith == "") {
		bWaitingForSinglePlayer = false;
	}
	if (pPlayer.DuoOffer != "") {
		delete(MapDuoOffers, pPlayer.DuoOffer);
		pPlayer.DuoOffer = "";
	}
	SetLastUpdated();
}

func Leave(pPlayer *players.EntPlayer, bGameStart bool) { //Players must be locked outside
	iPlayer := FindPlayerInArray(pPlayer, arQueue);
	if (iPlayer != -1) {
		arQueue = append(arQueue[:iPlayer], arQueue[iPlayer+1:]...);
		IPlayersCount--;
		if (pPlayer.IsReadyUpRequested && pPlayer.IsReadyConfirmed) {
			IReadyPlayers--;
		}
		if (bWaitingForSinglePlayer && pPlayer.DuoWith == "") {
			bWaitingForSinglePlayer = false;
		}

		pPlayer.NextQueueingAllowed = time.Now().UnixMilli() + 500; //500ms delay
		if (!bGameStart) {
			pPlayer.DuoWith = "";
		}

		pPlayer.IsInQueue = false;
		pPlayer.InQueueSince = 0;
		pPlayer.IsReadyUpRequested = false;
		pPlayer.IsReadyConfirmed = false;

		if (pPlayer == PLongestWaitPlayer) {
			if (IPlayersCount == 0) {
				PLongestWaitPlayer = nil;
			} else {
				PLongestWaitPlayer = GetLongestWaitPlayer();
			}
		}

		SetLastUpdated();
	}
}

func ReadyUp(pPlayer *players.EntPlayer) { //Players must be locked outside
	if (!pPlayer.IsReadyConfirmed) {
		pPlayer.IsReadyConfirmed = true;
		IReadyPlayers++;
	}
	SetLastUpdated();
}

func RequestReadyUp() { //queue is >= 8 players
	BIsInReadyUp = true;
	i64InReadyUpSince = time.Now().UnixMilli();
	pPlayerReadyUpReason = PLongestWaitPlayer;
	IReadyPlayers = 0;
	for _, pPlayer := range arQueue {
		pPlayer.IsReadyUpRequested = true;
		pPlayer.IsReadyConfirmed = false;
	}
	SetLastUpdated();
}

func StopReadyUp() {
	BIsInReadyUp = false;
	i64InReadyUpSince = 0;
	IReadyPlayers = 0;
	pPlayerReadyUpReason = nil;
	for _, pPlayer := range arQueue {
		pPlayer.IsReadyUpRequested = false;
		pPlayer.IsReadyConfirmed = false;
	}
}

func KickUnready() {
	var arKickPlayers []*players.EntPlayer;
	for _, pPlayer := range arQueue {
		if (pPlayer.IsReadyUpRequested && !pPlayer.IsReadyConfirmed) {
			arKickPlayers = append(arKickPlayers, pPlayer);
			if (pPlayer.DuoWith != "") {
				pDuoPlayer, bFound := players.MapPlayers[pPlayer.DuoWith];
				if (bFound && FindPlayerInArray(pDuoPlayer, arKickPlayers) == -1) {
					arKickPlayers = append(arKickPlayers, pDuoPlayer);
				}
			}
		}
	}
	for _, pPlayer := range arKickPlayers {
		Leave(pPlayer, false);
	}
}

func KickOffline() {
	var arKickPlayers []*players.EntPlayer;
	for _, pPlayer := range arQueue {
		if (!pPlayer.IsOnline) {
			arKickPlayers = append(arKickPlayers, pPlayer);
			if (pPlayer.DuoWith != "") {
				pDuoPlayer, bFound := players.MapPlayers[pPlayer.DuoWith];
				if (bFound && FindPlayerInArray(pDuoPlayer, arKickPlayers) == -1) {
					arKickPlayers = append(arKickPlayers, pDuoPlayer);
				}
			}
		}
	}
	for _, pPlayer := range arKickPlayers {
		Leave(pPlayer, false);
	}
}