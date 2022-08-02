package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomString(n int, letters string) (string, error) {
	ret := make([]byte, n);
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))));
		if err != nil {
			return "", err;
		}
		ret[i] = letters[num.Int64()];
	}
	return string(ret), nil;
}

func MaxValInt64(val1 int64, val2 int64) int64 {
	if (val1 > val2) {
		return val1;
	}
	return val2;
}

func GetStringIdxInArray(sValueBuffer string, arBuffer []string) int {
	iMax := len(arBuffer);
	for i := 0; i < iMax; i++ {
		if (arBuffer[i] == sValueBuffer) {
			return i;
		}
	}
	return -1;
}
