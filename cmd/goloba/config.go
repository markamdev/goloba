package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	envNamePort    = "GOLOBA_PORT"
	envNameTargets = "GOLOBA_TARGETS"
)

type config struct {
	Port    uint     `json:"port"`
	Servers []string `json:"servers"`
}

func loadConfig(fileName string) (config, error) {
	envConfig, err := loadEnvs()
	if err != nil {
		return config{}, err
	}
	if envConfig.Port > 0 && len(envConfig.Servers) > 0 {
		// necessary config given using environment variables
		return envConfig, nil
	}
	// not all data given by environment variables - has to load file
	var fileConfig config
	confData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return config{}, errors.New("Failed to open configuration file: " + fileName)
	}
	if err := json.Unmarshal(confData, &fileConfig); err != nil {
		return config{}, errors.New("Failed to process configuration in file: " + err.Error())
	}
	if envConfig.Port > 0 {
		fileConfig.Port = envConfig.Port
	}
	if len(envConfig.Servers) > 0 {
		fileConfig.Servers = envConfig.Servers
	}

	return fileConfig, nil
}

func loadEnvs() (config, error) {
	var result config
	// ex. GOLOBA_PORT=8000
	envPort := os.Getenv(envNamePort)
	if len(envPort) > 0 {
		val, err := strconv.ParseUint(envPort, 10, 32)
		if err != nil {
			return result, errors.New("Invalid PORT variable")
		}
		result.Port = uint(val)
	}
	// ex. GOLOBA_TARGETS="127.0.0.1:9000;10.0.0.1:9000"
	envTargets := os.Getenv(envNameTargets)
	if len(envTargets) > 0 {
		parts := strings.Split(envTargets, ";")
		result.Servers = parts
	}

	return result, nil
}
