package main

import (
	"runtime/pprof"
	"os"
	"runtime"
	"time"
	//"net/http"
	//_ "net/http/pprof"
)


func WatchMemory() {
	//go http.ListenAndServe(":24522", nil);
	for {
		time.Sleep(1 * time.Second);

		var memStat runtime.MemStats;
		runtime.ReadMemStats(&memStat);

		if (memStat.Sys > 5 * 1024 * 1024 * 1024) { //more than 5 gigs
			hMemStatDump, errMemStatDump := os.OpenFile("/root/memleakdump", os.O_CREATE|os.O_WRONLY, 0777);
			if (errMemStatDump == nil) {
				pprof.WriteHeapProfile(hMemStatDump);
				hMemStatDump.Close();
				return;
			}
		}
	}
}
