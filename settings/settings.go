package settings

import (
	"fmt"
	"flag"
	"github.com/buger/jsonparser"
	"io/ioutil"
)

var FilePath string;
var ListenPort string;

func Parse() bool {
	CommandLine();
	return ConfigFile();
}

func CommandLine() {
	oFilePath := flag.String("-config-path", "./settings.json", "Path to the settings.json file");

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

	ListenPort, errError = jsonparser.GetString(byData, "port");
	if (errError != nil) {
		fmt.Printf("Error reading config file: %s\n", errError);
		return false;
	}

	return true;
}
