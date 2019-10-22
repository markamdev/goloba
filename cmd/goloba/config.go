package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type config struct {
	Port    uint     `json:"port"`
	Servers []string `json:"servers"`
}

func loadConfigFile(fileName string) (config, error) {

	var result config
	confData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return config{}, errors.New("Failed to open configuration file: " + fileName)
	}
	if err := json.Unmarshal(confData, &result); err != nil {
		return config{}, errors.New("Failed to process configuration in file: " + err.Error())
	}

	return result, nil
}
