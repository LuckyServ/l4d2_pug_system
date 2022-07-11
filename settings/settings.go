package settings

import (
	"fmt"
	"flag"
	"github.com/buger/jsonparser"
	"io/ioutil"
)

var FilePath string;
var ListenPort string;
var BackendAuthKey string;
var DefaultMmrUncertainty float32;
var MmrStable float32;
var HomeDomain string;
var BackendDomain string;
var BrokenMode bool;

//database settings
var DatabaseHost string;
var DatabasePort string;
var DatabaseUsername string;
var DatabasePassword string;
var DatabaseName string;


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

	return true;
}
