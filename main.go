package main

import (
	"fmt"
	"./settings"
	"./players"
	"./database"
	"./lobby"
	"./api"
	"./players/auth"
	"time"
	//"./utils"
	//"crypto/rand"
	//"math/big"
	//"encoding/json"
)


func main() {
	fmt.Printf("Started.\n");

	//Parse settings
	if (!settings.Parse()) {
		return;
	}

	//Permanent global database connection
	if (!database.DatabaseConnect()) {
		return;
	}

	//Restore from database
	if (!players.RestorePlayers()) {
		return;
	}
	if (!auth.RestoreSessions()) {
		return;
	}
	lobby.I64LastLobbyListUpdate = time.Now().UnixMilli();

	//HTTP server init
	api.GinInit();

	go players.WatchOnline();






	//Simulation
	//var arMmrs []int = []int{1335,1382,1398,1875,850,2525,1497,2119,793,1526,1151,2096,1744,964,1443,2415,388,709,1698,1633,2014,1879,2130,2068,1995,1684,2838,1163,1995,1678,2545,2877,1845,732,1375,388,1436,914,1099,920,1811,1714,1431,876,819,2127,1853,1713,2182,1763,964,1690,2392,1757,1993,2392,2258,1824,1212,659,2168,1870,1077,1554,2196,1160,3232,3221,1111,2436,2940,2055,1831,1790,1958,2316,1116,1055,247,1081,2996,2618,1818,655,3404,1704,3526,1372,1675,2096,1130,1770,1594,1134,1358,1596,2404,1581,2546,2095,1986,1053,2358,1597,1756,2159,2463,1651,1029,1677,1320,792,1379,1500,2151,2203,903,2194,2337,1473,1864,991,1815,934,1504,1982,1975,2076,3265,1996,1737,1898,955,1600,1739};
	//var arMmrs []int = []int{1398,860,899,2200,1409,1879,2017,1073,1494,1375,388,1154,1894,825,1799,2192,1381,2229,1244,2940,1958,1596,2618,1818,3526,965,1675,1986,1898,414,1757,1756,1830,1500,934,2337,1511,541,949,2907,2860,1255,1119,1872,1961,1882,3030,1104,1471,2454,1404,2904,1719,2955,2992,1758,722,533,1597,1999,1424,796,1756,1124,2529,1454,1281,1};
	/*for i := 0; i < len(arMmrs); i++ {
		sRandSteamID, _ := utils.GenerateRandomString(17, "123456789");
		sName, _ := utils.GenerateRandomString(10, "abcdefghijklmnopqrstuvwxyz");
		//num, _ := rand.Int(rand.Reader, big.NewInt(int64(3000)));
		pPlayer := &players.EntPlayer{
			SteamID64:			sRandSteamID,
			Mmr:				arMmrs[i],//int(num.Int64() + 300),
			MmrUncertainty:		33.3,
			NicknameBase64:		sName,
			ProfValidated:		true,
			RulesAccepted:		true,
			LastActivity:		2000000000000,
			IsOnline:			true,
		};
		players.MapPlayers[sRandSteamID] = pPlayer;
		players.ArrayPlayers = append(players.ArrayPlayers, pPlayer);
		oSession := auth.EntSession{
			SteamID64:	sRandSteamID,
			Since:		time.Now().UnixMilli(),
		};
		sAuthBuffer, _ := utils.GenerateRandomString(32, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz");
		auth.MapSessions[sAuthBuffer] = oSession;
	}

	iLen := len(players.ArrayPlayers);
	for i := 0; i < iLen; i++ {



		pPlayer := players.ArrayPlayers[i];

		arLobbies := lobby.GetJoinableLobbies(pPlayer.Mmr);
		iSize := len(arLobbies);
		if (iSize == 0) {

			if (lobby.Create(pPlayer)) {
			} else {
			}

		} else {
			//sort
			if (iSize > 1) {
				bSorted := false;
				for !bSorted {
					bSorted = true;
					for i := 1; i < iSize; i++ {
						if (arLobbies[i].CreatedAt < arLobbies[i - 1].CreatedAt) {
							arLobbies[i], arLobbies[i - 1] = arLobbies[i - 1], arLobbies[i]; //switch
							bSorted = false;
						}
					}
				}
			}
			sLobbyID := arLobbies[0].ID;

			if (lobby.Join(pPlayer, sLobbyID)) {
			} else {
			}

		}

	}
	empData, _ := json.Marshal(lobby.MapLobbies);
	fmt.Printf("%s\n", empData);*/





	//Block until shutdown command is received
	fmt.Printf("End: %v\n", <-api.ChShutdown);
}
