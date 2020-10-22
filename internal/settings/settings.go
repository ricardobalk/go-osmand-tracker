package settings

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// Config contains the possible configuration parametes that are available in the settings.json file.
type Config struct {
	Port  uint `json:"port"`
	Debug bool `json:"debug"`
}

var (
	errInvalidPort = errors.New("Invalid port")
)

// Read reads the configuration file, validates it and returns it.
func Read(filename string) (*Config, error) {
	settingsFile, err := os.Open(filename)

	if err != nil {
		homeDir, err := os.UserHomeDir()
		settingsFile, err = os.Open(homeDir + "/" + filename)

		if err != nil {
			return nil, err
		}
	}

	byteValue, err := ioutil.ReadAll(settingsFile)

	if err != nil {
		return nil, err
	}

	err = settingsFile.Close()

	if err != nil {
		return nil, err
	}

	var configFile Config

	err = json.Unmarshal(byteValue, &configFile)

	if err != nil {
		return nil, err
	}

	// Validate entered port
	if configFile.Port == 0 {
		return nil, errInvalidPort
	}

	if configFile.Debug {
		log.Printf("Successfully parsed settings file: %s\n", filename)
	}

	return &configFile, nil
}

// Write writes the configuration file and returns an error in case of failure.
func Write(filename string, config *Config) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, configBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// IsCorrupted checks the configuration file for corruption and/or invalid values, it returns true in case the settings file is corrupted.
func IsCorrupted(filename string) bool {
	b, err := Read(filename)

	// should be corrupted if empty or invalid port
	if b == nil || err == errInvalidPort {
		return true
	}

	return false
}
