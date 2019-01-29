package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type RSSRule struct {
	URL      string `json:"url"`
	Interval int    `json:"interval"`
}

type HTMLRule struct {
	URL                 string `json:"url"`
	Interval            int    `json:"interval"`
	ArticleSelector     string `json:"article_selector"`
	TitleSelector       string `json:"title_selector"`
	DescriptionSelector string `json:"description_selector"`
}

type Config struct {
	RSSRules  []RSSRule  `json:"rss"`
	HTMLRules []HTMLRule `json:"html"`
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
