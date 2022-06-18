package main

import (
)

func SetShutDown() (bool, int) {
	if (bStateShutdown) {
		return false, 2;
	}
	bStateShutdown = true;
	return true, 0;
}

func PerformShutDown() {
	chShutdown <- true;
}
