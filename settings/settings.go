package settings

import (
	"fmt"
	"flag"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"sort"
	"strings"
	"strconv"
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

var MapConfoglConfigs map[int]string = make(map[int]string);
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


	//Confogl configs section
	sConfoglMmrs, errConfoglMmrs := jsonparser.GetString(byData, "lobby", "confogl_configs", "mmr");
	sConfoglConfigs, errConfoglConfigs := jsonparser.GetString(byData, "lobby", "confogl_configs", "config");
	if (errConfoglMmrs != nil || errConfoglConfigs != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	arConfoglMmrs := strings.Split(sConfoglMmrs, ",");
	arConfoglConfigs := strings.Split(sConfoglConfigs, ",");
	if (len(arConfoglMmrs) != len(arConfoglConfigs) || len(arConfoglMmrs) == 0 || len(arConfoglConfigs) == 0) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}
	for i, _ := range arConfoglMmrs {
		iConfMmr, errConfMmr := strconv.Atoi(arConfoglMmrs[i]);
		sConf := arConfoglConfigs[i];
		if (errConfMmr != nil || iConfMmr <= 0 || sConf == "") {
			fmt.Printf("Error reading config file: %s\n", errError);
			return false;
		}
		MapConfoglConfigs[iConfMmr] = sConf;
		ArrayConfoglConfigsMmrs = append(ArrayConfoglConfigsMmrs, iConfMmr);
	}
	sort.Ints(ArrayConfoglConfigsMmrs);


	//Map pool section
	bErrorReadingMapPool := false;
	jsonparser.ArrayEach(byData, func(valueCampaign []byte, dataType jsonparser.ValueType, offset int, err error) {
		sCampaignName, _ := jsonparser.GetString(valueCampaign, "name");
		var arCampaign []string;
		jsonparser.ArrayEach(valueCampaign, func(valueMap []byte, dataType jsonparser.ValueType, offset int, err error) {
			sMap := string(valueMap);
			if (sMap != "") {
				arCampaign = append(arCampaign, sMap);
			} else {
				bErrorReadingMapPool = true;
			}
		}, "maps");
		if (len(arCampaign) > 0 && sCampaignName != "") {
			MapPool = append(MapPool, arCampaign);
			CampaignNames = append(CampaignNames, sCampaignName);
		} else {
			bErrorReadingMapPool = true;
		}
	}, "map_pool");
	if (bErrorReadingMapPool) {
		fmt.Printf("Error reading config file on map pool section\n");
		return false;
	}


	return true;
}
