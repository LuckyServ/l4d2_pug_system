package games

import (
	"../players"
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
	//Pair players
	//Ping servers
	//Select server
	//Game proceeds
	//Game ended, settle results
	//Destroy Game
}