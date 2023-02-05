package settings

import (
	"fmt"
	"flag"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"sort"
	//"../utils"
)

var FilePath string;
var LogPath string;
var ListenPort string;
var MapBackendAuthKeys map[string]bool;
var DefaultMmrUncertainty float32;
var IncreaseMmrUncertainty float32;
var MmrStable float32;
var MmrAbsoluteWin float32;
var MmrMinimumWin float32;
var MmrDiffGuaranteedWin float32;
var MmrGrades int;
var HomeDomain string;
var BackendDomain string;
var BrokenMode bool;

var DatabaseHost string;
var DatabasePort string;
var DatabaseUsername string;
var DatabasePassword string;
var DatabaseName string;

var SteamApiKey string;
var MinVersusGamesPlayed int;
var DefaultMaxMmr int;

var OnlineTimeout int64;
var ReadyUpTimeout int64;
var PlayerAuthExpire int64;

var AuthPerHour int;
var ExternalAPICooldown int64;

var GameServers []GameServer;

var MaxPingWait int;
var AvailGameSrvsMaxTries int;
var FirstReadyUpExpire int64;
var MaxAbsentSeconds int64;
var MaxSingleAbsentSeconds int64;
var MinPlayersCount int;

var RQDefaultPenalty int;

var ChatStoreMaxMsgs int;
var ChatMsgDelay int64;
var ChatMaxChars int;

var BanHistoryForgetIn int64;
var BanRQFirst int64;
var BanRQSecond int64;
var BanRQThird int64;
var BanRQReason string;

var BanListPagination int;

var SmurfHost string;
var SmurfAuthKey string;

var TwitchClientID string;
var TwitchSecret string;


type GameServer struct {
	IP				string		`json:"ip"`
	Domain			string		`json:"domain"`
	Ports			[]string	`json:"ports"`
	LowerPriority	int			`json:"lower_priority"`
}
var MapProxies = make(map[string]string);

type ConfoglConfig struct {
	CodeName		string
	Name			string
	MmrMax			int
}

var MapConfoglConfigs map[int]ConfoglConfig = make(map[int]ConfoglConfig);
var ArrayConfoglConfigsMmrs []int;


type Campaign struct {
	Name			string		`json:"name"`
	Maps			[]string	`json:"maps"`
	DownloadLink	string		`json:"download"`
	LastUpdated		int64		`json:"last_updated"`
}
var MapPool []Campaign;
var NewestCustomMap int64;



func Parse() bool {
	CommandLine();
	return ConfigFile();
}

func CommandLine() {
	oFilePath := flag.String("config-path", "./settings.json", "Path to the settings.json file");
	oLogPath := flag.String("log-path", "./pugs.log", "Path to the log file");

	flag.Parse();

	FilePath = *oFilePath;
	LogPath = *oLogPath;
}

func ConfigFile() bool {
	byData, errFile := ioutil.ReadFile(FilePath);
	if (errFile != nil) {
		fmt.Printf("Error reading config file: %s\n", errFile);
		return false;
	}
	var errError error;
	var f64Buffer float64;
	var i64Buffer int64;

	ListenPort, errError = jsonparser.GetString(byData, "listen_port");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	f64Buffer, errError = jsonparser.GetFloat(byData, "mmr", "default_uncertainty");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	DefaultMmrUncertainty = float32(f64Buffer);

	f64Buffer, errError = jsonparser.GetFloat(byData, "mmr", "increase_daily");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	IncreaseMmrUncertainty = float32(f64Buffer);

	f64Buffer, errError = jsonparser.GetFloat(byData, "mmr", "consider_stable");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	MmrStable = float32(f64Buffer);

	f64Buffer, errError = jsonparser.GetFloat(byData, "mmr", "absolute_win");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	MmrAbsoluteWin = float32(f64Buffer);

	f64Buffer, errError = jsonparser.GetFloat(byData, "mmr", "minimum_win");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	MmrMinimumWin = float32(f64Buffer);

	f64Buffer, errError = jsonparser.GetFloat(byData, "mmr", "mmr_diff_guaranteed_win");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	MmrDiffGuaranteedWin = float32(f64Buffer);

	HomeDomain, errError = jsonparser.GetString(byData, "domain", "home");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	BackendDomain, errError = jsonparser.GetString(byData, "domain", "backend");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	BrokenMode, errError = jsonparser.GetBoolean(byData, "broken_state");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	DatabaseHost, errError = jsonparser.GetString(byData, "database", "host");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	DatabasePort, errError = jsonparser.GetString(byData, "database", "port");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	DatabaseUsername, errError = jsonparser.GetString(byData, "database", "user");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	DatabasePassword, errError = jsonparser.GetString(byData, "database", "password");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	DatabaseName, errError = jsonparser.GetString(byData, "database", "name");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	SteamApiKey, errError = jsonparser.GetString(byData, "steam_api_key");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	i64Buffer, errError = jsonparser.GetInt(byData, "min_versus_games_played");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	MinVersusGamesPlayed = int(i64Buffer);

	i64Buffer, errError = jsonparser.GetInt(byData, "mmr", "default_max");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	DefaultMaxMmr = int(i64Buffer);

	i64Buffer, errError = jsonparser.GetInt(byData, "mmr", "mmr_grades");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	MmrGrades = int(i64Buffer);

	i64Buffer, errError = jsonparser.GetInt(byData, "timeouts", "online_seconds");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	OnlineTimeout = i64Buffer * 1000;

	i64Buffer, errError = jsonparser.GetInt(byData, "ratelimits", "auth_per_hour");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	AuthPerHour = int(i64Buffer);

	i64Buffer, errError = jsonparser.GetInt(byData, "ratelimits", "external_api_cooldown_sec");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	ExternalAPICooldown = i64Buffer * 1000;

	i64Buffer, errError = jsonparser.GetInt(byData, "timeouts", "readyup_seconds");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	ReadyUpTimeout = i64Buffer * 1000;

	i64Buffer, errError = jsonparser.GetInt(byData, "timeouts", "player_auth_expire_days");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	PlayerAuthExpire = i64Buffer * 24 * 60 * 60 * 1000;

	i64Buffer, errError = jsonparser.GetInt(byData, "game", "max_wait_pings_seconds");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	MaxPingWait = int(i64Buffer);

	i64Buffer, errError = jsonparser.GetInt(byData, "game", "get_available_servers_max_tries");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	AvailGameSrvsMaxTries = int(i64Buffer);

	i64Buffer, errError = jsonparser.GetInt(byData, "game", "new_game_first_readyup_minutes");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	FirstReadyUpExpire = i64Buffer * 60;

	i64Buffer, errError = jsonparser.GetInt(byData, "game", "max_absent_minutes");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	MaxAbsentSeconds = i64Buffer * 60;

	i64Buffer, errError = jsonparser.GetInt(byData, "game", "max_single_absent_minutes");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	MaxSingleAbsentSeconds = i64Buffer * 60;

	i64Buffer, errError = jsonparser.GetInt(byData, "game", "game_ends_if_players_disconnected");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	MinPlayersCount = 8 - int(i64Buffer);

	i64Buffer, errError = jsonparser.GetInt(byData, "game", "ragequit_penalty", "default");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	RQDefaultPenalty = int(i64Buffer);

	i64Buffer, errError = jsonparser.GetInt(byData, "chat", "store_n_messages");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	ChatStoreMaxMsgs = int(i64Buffer);

	i64Buffer, errError = jsonparser.GetInt(byData, "chat", "message_delay_ms");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	ChatMsgDelay = i64Buffer;

	i64Buffer, errError = jsonparser.GetInt(byData, "chat", "max_chars");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	ChatMaxChars = int(i64Buffer);

	i64Buffer, errError = jsonparser.GetInt(byData, "ragequit_bans", "forget_ban_history_in_days");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	BanHistoryForgetIn = i64Buffer * 24 * 60 * 60 * 1000;

	i64Buffer, errError = jsonparser.GetInt(byData, "ragequit_bans", "first_ban_hours");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	BanRQFirst = i64Buffer * 60 * 60 * 1000;

	i64Buffer, errError = jsonparser.GetInt(byData, "ragequit_bans", "second_ban_hours");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	BanRQSecond = i64Buffer * 60 * 60 * 1000;

	i64Buffer, errError = jsonparser.GetInt(byData, "ragequit_bans", "third_ban_hours");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	BanRQThird = i64Buffer * 60 * 60 * 1000;

	BanRQReason, errError = jsonparser.GetString(byData, "ragequit_bans", "reason_text");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	i64Buffer, errError = jsonparser.GetInt(byData, "banlist_items_per_page");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	BanListPagination = int(i64Buffer);

	SmurfHost, errError = jsonparser.GetString(byData, "smurf_detector", "host");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	SmurfAuthKey, errError = jsonparser.GetString(byData, "smurf_detector", "auth_key");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	TwitchClientID, errError = jsonparser.GetString(byData, "twitch", "client_id");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	TwitchSecret, errError = jsonparser.GetString(byData, "twitch", "client_secret");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}


	//Confogl configs section
	bErrorReadingConfoglConfigs := true;
	jsonparser.ArrayEach(byData, func(valueConfoglConfig []byte, dataType jsonparser.ValueType, offset int, err error) {


		sConfoglConfigName, errConfoglConfigName := jsonparser.GetString(valueConfoglConfig, "name");
		sConfoglConfigCodename, errConfoglConfigCodename := jsonparser.GetString(valueConfoglConfig, "codename");
		i64ConfoglConfigMaxMmr, errConfoglConfigMaxMmr := jsonparser.GetInt(valueConfoglConfig, "max_mmr");

		if (errConfoglConfigName == nil && errConfoglConfigCodename == nil && errConfoglConfigMaxMmr == nil && sConfoglConfigName != "" && sConfoglConfigCodename != "" && i64ConfoglConfigMaxMmr != 0) {
			iConfMmr := int(i64ConfoglConfigMaxMmr);

			oConfoglConf := ConfoglConfig{
				CodeName:		sConfoglConfigCodename,
				Name:			sConfoglConfigName,
				MmrMax:			iConfMmr,
			};

			MapConfoglConfigs[iConfMmr] = oConfoglConf;
			ArrayConfoglConfigsMmrs = append(ArrayConfoglConfigsMmrs, iConfMmr);

			bErrorReadingConfoglConfigs = false;
		}




	}, "confogl_configs");
	if (bErrorReadingConfoglConfigs) {
		fmt.Printf("Error reading config file on Confogl configs list\n");
		return false;
	}
	sort.Ints(ArrayConfoglConfigsMmrs);


	//Map pool section
	if (!UpdateMapsFromJSON(byData)) {
		return false;
	}

	//Gameservers section
	if (!UpdateServersFromJSON(byData)) {
		return false;
	}
	
	//Backend auth keys
	if (!UpdateAuthKeysFromJSON(byData)) {
		return false;
	}

	return true;
}

func RefreshAuthKeys() {
	byData, errFile := ioutil.ReadFile(FilePath);
	if (errFile != nil) {
		fmt.Printf("Error reading config file: %s\n", errFile);
		return;
	}
	UpdateAuthKeysFromJSON(byData);
}

func UpdateAuthKeysFromJSON(byData []byte) bool {
	bErrorReadingAuthKeys := true;
	MapBackendAuthKeys = make(map[string]bool, 0);
	jsonparser.ArrayEach(byData, func(valueAuthKey []byte, dataType jsonparser.ValueType, offset int, err error) {
		sAuthKey := string(valueAuthKey);
		if (sAuthKey != "") {
			MapBackendAuthKeys[sAuthKey] = true;
			bErrorReadingAuthKeys = false;
		}
	}, "backend_auth_keys");
	if (bErrorReadingAuthKeys) {
		fmt.Printf("Error reading config file on auth keys section\n");
		return false;
	}
	return true;
}

func RefreshMaps() {
	byData, errFile := ioutil.ReadFile(FilePath);
	if (errFile != nil) {
		fmt.Printf("Error reading config file: %s\n", errFile);
		return;
	}
	UpdateMapsFromJSON(byData);
}

func UpdateMapsFromJSON(byData []byte) bool {
	bErrorReadingMapPool := true;
	NewestCustomMap = 0;
	MapPool = make([]Campaign, 0);
	jsonparser.ArrayEach(byData, func(valueCampaign []byte, dataType jsonparser.ValueType, offset int, err error) {
		sCampaignName, _ := jsonparser.GetString(valueCampaign, "name");
		sDownloadLink, _ := jsonparser.GetString(valueCampaign, "download");
		i64UpdatedAt, _ := jsonparser.GetInt(valueCampaign, "last_updated");
		var arMaps []string;
		jsonparser.ArrayEach(valueCampaign, func(valueMap []byte, dataType jsonparser.ValueType, offset int, err error) {
			sMap := string(valueMap);
			if (sMap != "") {
				arMaps = append(arMaps, sMap);
				bErrorReadingMapPool = false;
			}
		}, "maps");
		if (len(arMaps) > 0 && sCampaignName != "" && i64UpdatedAt > 0) {

			oCampaign := Campaign{
				Name:			sCampaignName,
				Maps:			arMaps,
				DownloadLink:	sDownloadLink,
				LastUpdated:	i64UpdatedAt,
			}

			MapPool = append(MapPool, oCampaign);
			bErrorReadingMapPool = false;
			if (i64UpdatedAt > NewestCustomMap) {
				NewestCustomMap = i64UpdatedAt;
			}
		}
	}, "map_pool");
	if (bErrorReadingMapPool) {
		fmt.Printf("Error reading config file on map pool section\n");
		return false;
	}
	return true;
}

func RefreshServers() {
	byData, errFile := ioutil.ReadFile(FilePath);
	if (errFile != nil) {
		fmt.Printf("Error reading config file: %s\n", errFile);
		return;
	}
	UpdateServersFromJSON(byData);
}

func UpdateServersFromJSON(byData []byte) bool {
	bErrorReadingGameServers := true;
	GameServers = make([]GameServer, 0);
	MapProxies = make(map[string]string);
	jsonparser.ArrayEach(byData, func(valueServer []byte, dataType jsonparser.ValueType, offset int, err error) {

		sDomain, _ := jsonparser.GetString(valueServer, "domain");
		sServIP, _ := jsonparser.GetString(valueServer, "ip");
		sProxyIP, _ := jsonparser.GetString(valueServer, "proxy");
		i64LowerPriority, _ := jsonparser.GetInt(valueServer, "lower_priority");
		iLowerPriority := int(i64LowerPriority);
		oGameServer := GameServer{
			IP:				sServIP,
			Domain:			sDomain,
			LowerPriority:	iLowerPriority,
		};
		
		jsonparser.ArrayEach(valueServer, func(valuePort []byte, dataType jsonparser.ValueType, offset int, err error) {
			sPORT := string(valuePort);
			if (sPORT != "") {
				oGameServer.Ports = append(oGameServer.Ports, sPORT);
				if (sProxyIP != "") {
					MapProxies[sServIP + ":" + sPORT] = sProxyIP + ":" + sPORT;
				}
				if (bErrorReadingGameServers == true) {
					bErrorReadingGameServers = false;
				}
			}
		}, "ports");

		GameServers = append(GameServers, oGameServer);

	}, "game_servers");
	if (bErrorReadingGameServers) {
		fmt.Printf("Error reading config file on gameservers section\n");
		return false;
	}
	return true;
}
