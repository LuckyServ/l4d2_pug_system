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
var ListenPort string;
var BackendAuthKey string;
var DefaultMmrUncertainty float32;
var MmrStable float32;
var HomeDomain string;
var BackendDomain string;
var BrokenMode bool;
var NoNewLobbies bool;

var DatabaseHost string;
var DatabasePort string;
var DatabaseUsername string;
var DatabasePassword string;
var DatabaseName string;

var SteamApiKey string;
var MinVersusGamesPlayed int;
var DefaultMaxMmr int;
var OnlineMmrRange int;

var IdleTimeout int64;
var OnlineTimeout int64;
var ReadyUpTimeout int64;
var LobbyFillTimeout int64;

var JoinLobbyCooldown int64;
var AuthPerHour int;
var ProfValidateCooldown int64;

var GameServers []string;
var HardwareServers []string; //parallel with HardwareDomains
var HardwareDomains []string; //parallel with HardwareServers

var MaxPingWait int;
var AvailGameSrvsMaxTries int;
var FirstReadyUpExpire int64;
var MaxAbsentSeconds int64;


type ConfoglConfig struct {
	CodeName		string
	Name			string
	MmrMax			int
}

var MapConfoglConfigs map[int]ConfoglConfig = make(map[int]ConfoglConfig);
var ArrayConfoglConfigsMmrs []int;

var MapPool [][]string;
var CampaignNames []string; //parallel with MapPool



func Parse() bool {
	CommandLine();
	return ConfigFile();
}

func CommandLine() {
	oFilePath := flag.String("config-path", "./settings.json", "Path to the settings.json file");

	flag.Parse();

	FilePath = *oFilePath;
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

	BackendAuthKey, errError = jsonparser.GetString(byData, "backend_auth_key");
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

	f64Buffer, errError = jsonparser.GetFloat(byData, "mmr", "consider_stable");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	MmrStable = float32(f64Buffer);

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

	NoNewLobbies, errError = jsonparser.GetBoolean(byData, "no_new_lobbies");
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

	i64Buffer, errError = jsonparser.GetInt(byData, "lobby", "online_mmr_range");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	OnlineMmrRange = int(i64Buffer);

	i64Buffer, errError = jsonparser.GetInt(byData, "timeouts", "online_minutes");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	OnlineTimeout = i64Buffer * 60 * 1000;

	i64Buffer, errError = jsonparser.GetInt(byData, "timeouts", "idle_minutes");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	IdleTimeout = i64Buffer * 60 * 1000;

	i64Buffer, errError = jsonparser.GetInt(byData, "ratelimits", "lobby_join_cooldown_sec");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	JoinLobbyCooldown = i64Buffer * 1000;

	i64Buffer, errError = jsonparser.GetInt(byData, "ratelimits", "auth_per_hour");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	AuthPerHour = int(i64Buffer);

	i64Buffer, errError = jsonparser.GetInt(byData, "ratelimits", "prof_validate_cooldown_sec");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	ProfValidateCooldown = i64Buffer * 1000;

	i64Buffer, errError = jsonparser.GetInt(byData, "timeouts", "readyup_seconds");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	ReadyUpTimeout = i64Buffer * 1000;

	i64Buffer, errError = jsonparser.GetInt(byData, "timeouts", "lobby_fill_minutes");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	LobbyFillTimeout = i64Buffer * 60 * 1000;

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




	}, "lobby", "confogl_configs");
	if (bErrorReadingConfoglConfigs) {
		fmt.Printf("Error reading config file on Confogl configs list\n");
		return false;
	}
	sort.Ints(ArrayConfoglConfigsMmrs);


	//Map pool section
	bErrorReadingMapPool := true;
	jsonparser.ArrayEach(byData, func(valueCampaign []byte, dataType jsonparser.ValueType, offset int, err error) {
		sCampaignName, _ := jsonparser.GetString(valueCampaign, "name");
		var arCampaign []string;
		jsonparser.ArrayEach(valueCampaign, func(valueMap []byte, dataType jsonparser.ValueType, offset int, err error) {
			sMap := string(valueMap);
			if (sMap != "") {
				arCampaign = append(arCampaign, sMap);
				bErrorReadingMapPool = false;
			}
		}, "maps");
		if (len(arCampaign) > 0 && sCampaignName != "") {
			MapPool = append(MapPool, arCampaign);
			CampaignNames = append(CampaignNames, sCampaignName);
			bErrorReadingMapPool = false;
		}
	}, "map_pool");
	if (bErrorReadingMapPool) {
		fmt.Printf("Error reading config file on map pool section\n");
		return false;
	}


	//Gameservers section
	bErrorReadingGameServers := true;
	jsonparser.ArrayEach(byData, func(valueServer []byte, dataType jsonparser.ValueType, offset int, err error) {

		sDomain, _ := jsonparser.GetString(valueServer, "domain");
		sServIP, _ := jsonparser.GetString(valueServer, "ip");
		if (sDomain != "" && sServIP != "") {
			HardwareDomains = append(HardwareDomains, sDomain);
			HardwareServers = append(HardwareServers, sServIP);
		}
		jsonparser.ArrayEach(valueServer, func(valuePort []byte, dataType jsonparser.ValueType, offset int, err error) {
			sPORT := string(valuePort);
			if (sPORT != "") {
				GameServers = append(GameServers, sServIP + ":" + sPORT);
				bErrorReadingGameServers = false;
			}
		}, "ports");

	}, "game_servers");
	if (bErrorReadingGameServers) {
		fmt.Printf("Error reading config file on gameservers section\n");
		return false;
	}


	return true;
}
