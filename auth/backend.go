package auth

import (
	"../settings"
)

func Backend(sKey string) bool {
	if (sKey == settings.BackendAuthKey) {
		return true;
	}
	return false;
}
