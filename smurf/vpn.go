package smurf

import (
)

type EntVPNInfo struct {
	IsVPN		bool
	UpdatedAt	int64
}

var mapIPs = make(map[string]EntVPNInfo);
