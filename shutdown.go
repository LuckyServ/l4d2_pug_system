package main

import (
)

func SetShutDown() (bool, string) {
	if (bStateShutdown) {
		return false, "Already shutting down";
	}
	bStateShutdown = true;
	return true, "";
}

func PerformShutDown() {
	chShutdown <- true;
}
