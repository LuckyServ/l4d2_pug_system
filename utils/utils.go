package utils

import (
	"crypto/rand"
	"math/big"
	"bytes"
	"time"
	"fmt"
)

var ChanUniqueString = make(chan string);

func Watchers() {
	for {
		select {
		case ChanUniqueString <- func()(string) {
			sRand, _ := GenerateRandomString(32, "0123456789abcdefghijklmnopqrstuvwxyz");
			return fmt.Sprintf("%d%s", time.Now().UnixNano(), sRand);
		}():
		}
		time.Sleep(1 * time.Nanosecond);
	}
}

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

func GetRandInt(iMin int, iMax int) (int, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(int64((iMax - iMin) + 1)));
	if err != nil {
		return 0, err;
	}
	return (iMin + int(num.Int64())), nil;
}

func MaxValInt64(val1 int64, val2 int64) int64 {
	if (val1 > val2) {
		return val1;
	}
	return val2;
}

func MaxValInt(val1 int, val2 int) int {
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

func GetIntIdxInArray(iValueBuffer int, arBuffer []int) int {
	iMax := len(arBuffer);
	for i := 0; i < iMax; i++ {
		if (arBuffer[i] == iValueBuffer) {
			return i;
		}
	}
	return -1;
}

func RemoveIntFromArray(iIndex int, arBuffer []int) []int {
	arBuffer[iIndex] = arBuffer[len(arBuffer)-1];
	return arBuffer[:len(arBuffer)-1];
}

func InsertDots(s string, n int) string {
	var buffer bytes.Buffer;
	var n_1 = n - 1;
	var l_1 = len(s) - 1;
	for i,rune := range s {
		buffer.WriteRune(rune);
		if i % n == n_1 && i != l_1 {
			buffer.WriteRune('.');
		}
	}
	return buffer.String();
}
