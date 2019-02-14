package main

import (
	"database/sql"
	"log"
	"time"
)

type Aggregator struct {
	db     *sql.DB
	config *Config
}

func fetchNews(interval int, fetchFunc func()) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		// Можно в принципе обойтись без запуска этой функции в горутине
		// сделано было для того, чтобы не задерживать этот цикл в случае, если интервал меньше времени получения новостей с источника
		go fetchFunc()
		_ = <-ticker.C
	}
}

func (agr Aggregator) collectNewsAndSave(newsChan <-chan News) {
	for {
		news := <-newsChan
		err := agr.saveNews(&news)
		if err != nil {
			log.Println("error while saving news:", err)
		}
	}
}

func NewAggregator(dbPath, configPath string) (*Aggregator, error) {
	db, err := newDatabase(dbPath)

	if err != nil {
		return nil, err
	}

	config, err := parseConfig(configPath)

	if err != nil {
		return nil, err
	}

	agr := &Aggregator{
		db:     db,
		config: config,
	}

	return agr, nil
}

func (agr Aggregator) Run() {
	newsChan := make(chan News, 200)

	go agr.collectNewsAndSave(newsChan)

	for _, rule := range agr.config.RSSRules {
		go fetchNews(rule.Interval, func() {
			parseRSS(rule, newsChan)
		})
	}

	for _, rule := range agr.config.HTMLRules {
		go fetchNews(rule.Interval, func() {
			parseHTML(rule, newsChan)
		})
	}
}
