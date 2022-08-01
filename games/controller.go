package games

import (
	"../players"
	"../settings"
	"../rating"
	"time"
)


func Control(pGame *EntGame) {

	//Create
	MuGames.Lock();
	players.MuPlayers.Lock();

	Create(pGame);
	pGame.State = StateCreated;
	SetLastUpdated(pGame.PlayersUnpaired);

	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Choose maps
	i64CmpgnIdx := time.Now().UnixNano() % int64(len(settings.CampaignNames));
	MuGames.Lock();
	players.MuPlayers.Lock();
	pGame.CampaignName = settings.CampaignNames[i64CmpgnIdx];
	pGame.Maps = settings.MapPool[i64CmpgnIdx];
	pGame.State = StateCampaignChosen;
	SetLastUpdated(pGame.PlayersUnpaired);
	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Pair players
	MuGames.Lock();
	players.MuPlayers.Lock();

	pGame.PlayersA, pGame.PlayersB = rating.Pair(pGame.PlayersUnpaired);

	pGame.State = StateTeamsPicked;
	SetLastUpdated(pGame.PlayersUnpaired);
	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Request pings
	MuGames.Lock();
	players.MuPlayers.Lock();
	for _, pPlayer := range pGame.PlayersUnpaired {
		pPlayer.GameServersPinged = false;
	}
	pGame.State = StateWaitPings;
	SetLastUpdated(pGame.PlayersUnpaired);
	MuGames.Unlock();
	players.MuPlayers.Unlock();


	//Wait for pings
	






	//Select server
	//Game proceeds
	//Game ended, settle results
	//Destroy Game
}