package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Rule struct {
	URL      string `json:"url"`
	Interval int    `json:"interval"`
}

type Config struct {
	Rules []Rule `json:"rules"`
}

func parseConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	rawConfigData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(rawConfigData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
