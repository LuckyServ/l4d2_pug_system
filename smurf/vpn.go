package smurf

import (
)

type EntVPNInfo struct {
	IsVPN		bool
	UpdatedAt	int64
}

var mapIPs = make(map[string]EntVPNInfo);
var ChanAnnounce = make(chan string);
var ChanCheckVPN = make(chan bool);
var ChanClear = make(chan bool);
