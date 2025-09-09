package env

import (
	"encoding/json"
	cerr "github.com/jeanfrancoisgratton/customError"
	"os"
	"path/filepath"
	"strings"
)

// LoadEnvironmentFile : Load the JSON env file in the user's .config/JFG/vclt directory, and store it into a data type (struct)
func LoadEnvironmentFile() (Config_s, *cerr.CustomError) {
	var payload Config_s
	var ce *cerr.CustomError

	if !strings.HasSuffix(ConfigFile, ".json") {
		ConfigFile += ".json"
	}

	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vclt", ConfigFile)
	_, err := os.Stat(rcFile)
	// We need to create the environment file if it does not exist
	if os.IsNotExist(err) {
		payload, ce = getEnvParams(rcFile)
		if ce != nil {
			return Config_s{}, ce
		}
		es := &payload
		if ce = es.SaveEnvironmentFile(ConfigFile); ce != nil {
			return Config_s{}, ce
		}
	}

	jFile, err := os.ReadFile(rcFile)
	if err != nil {
		return Config_s{}, &cerr.CustomError{Title: "Error reading the file", Message: err.Error()}
	}
	err = json.Unmarshal(jFile, &payload)
	if err != nil {
		return Config_s{}, &cerr.CustomError{Title: "Error unmarshalling JSON", Message: err.Error()}
	} else {
		return payload, nil
	}
}

// SaveEnvironmentFile : Save the above structure into a JSON file in the user's .config/JFG/vclt directory
func (e *Config_s) SaveEnvironmentFile(outputfile string) *cerr.CustomError {
	if outputfile == "" {
		outputfile = ConfigFile
	}

	if !strings.HasSuffix(outputfile, ".json") {
		outputfile += ".json"
	}

	jStream, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return &cerr.CustomError{Title: "Error marshalling information", Message: err.Error()}
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vclt", outputfile)
	if err = os.WriteFile(rcFile, jStream, 0600); err != nil {
		return &cerr.CustomError{Title: "Error writing config file", Message: err.Error()}
	}

	return nil
}
