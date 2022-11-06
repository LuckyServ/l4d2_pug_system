package main

import (
	"runtime/pprof"
	"os"
	"runtime"
	"time"
	"net/http"
	_ "net/http/pprof"
)


/*func WatchMemory() {
	go http.ListenAndServe(":24522", nil);
	for {
		time.Sleep(100 * time.Millisecond);

		var memStat runtime.MemStats;
		runtime.ReadMemStats(&memStat);

		if (memStat.Sys > 2 * 1024 * 1024 * 1024) { //more than 2 gigs
			hMemStatDump, errMemStatDump := os.OpenFile("/root/memleakdump", os.O_CREATE|os.O_WRONLY, 0777);
			if (errMemStatDump == nil) {
				pprof.WriteHeapProfile(hMemStatDump);
				hMemStatDump.Close();
				return;
			}
		}
	}
}*/

func WatchDeadlocks() {
	runtime.SetBlockProfileRate(1);
	go http.ListenAndServe(":24522", nil);
	for {
		time.Sleep(100 * time.Millisecond);

		var memStat runtime.MemStats;
		runtime.ReadMemStats(&memStat);

		if (memStat.Sys > 2 * 1024 * 1024 * 1024) { //more than 2 gigs
			hMemStatDump, errMemStatDump := os.OpenFile("/root/blocksdump", os.O_CREATE|os.O_WRONLY, 0777);
			if (errMemStatDump == nil) {
				pprof.Lookup("block").WriteTo(hMemStatDump, 0);
				hMemStatDump.Close();
				return;
			}
		}
	}
}