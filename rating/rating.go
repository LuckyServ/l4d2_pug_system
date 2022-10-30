package rating

import (
	"github.com/robfig/cron/v3"
	"../players"
	"../settings"
	"../database"
)


func Watchers() {
	go SetCron();
}

func SetCron() {
	oCron := cron.New();

	oCron.AddFunc("0 7 * * *", func(){IncreaseUncertainty();});

	oCron.Run();
}

func IncreaseUncertainty() {
	players.MuPlayers.Lock();
	for _, pPlayer := range players.ArrayPlayers {
		pPlayer.MmrUncertainty = pPlayer.MmrUncertainty + settings.IncreaseMmrUncertainty;
		if (pPlayer.MmrUncertainty > settings.DefaultMmrUncertainty) {
			pPlayer.MmrUncertainty = settings.DefaultMmrUncertainty;
		}
	}
	players.MuPlayers.Unlock();
	go database.IncreaseUncertainty();
}