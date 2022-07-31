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

	MuGames.Unlock();
	players.MuPlayers.Unlock();

	//Choose maps
	i64CmpgnIdx := time.Now().UnixNano() % int64(len(settings.CampaignNames));
	MuGames.Lock();
	pGame.CampaignName = settings.CampaignNames[i64CmpgnIdx];
	pGame.Maps = settings.MapPool[i64CmpgnIdx];
	pGame.State = StateCampaignChosen;
	MuGames.Unlock();


	//Pair players
	MuGames.Lock();
	players.MuPlayers.Lock();
	
	pGame.PlayersA, pGame.PlayersB = rating.Pair(pGame.PlayersUnpaired);

	pGame.State = StateTeamsPicked;
	MuGames.Unlock();
	players.MuPlayers.Unlock();




	//Ping servers
	//Select server
	//Game proceeds
	//Game ended, settle results
	//Destroy Game
}