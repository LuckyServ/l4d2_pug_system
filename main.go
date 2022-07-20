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
	//var arMmrs []int = []int{1358,2321,1382,1398,1905,1751,747,1366,1021,1212,1691,1152,1134,1887,1826,2270,1358,1507,1473,2184,1600,3538,3274,2323,2309,1928,3188,2382,2794,2120,730,1081,1677,1632,827,2878,913,1706,1357,1500,1694,2555,1872,2153,1657,1011,1463,1669,1375,388,2512,2156,4,1707,915,782,1425,2127,695,2704,1409,1712,2841,1745,1993,1812,2488,409,1925,2003,515,3073,2449,1160,2947,1783,2353,1449,247,1081,2996,655,3404,3404,667,3502,2096,1963,1784,165,2285,1865,2212,1281,1990,2095,1986,2199,968,1883,1897,269,96,2184,2080,2229,2264,958,1753,1894,1320,1775,1363,1020,1500,901,3009,11,1923,2195,1478,915,2337,1221,756,1831,1815,1641,1049,1348,3230,1957,1885,0,4,1522,1382,2010,2347,2133,1232,27,2860,1943,1276,1594,1828,2351,1559,1784,1788,1799,2492,1119,1895,1719};
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
