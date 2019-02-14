package main

import (
	"database/sql"
	"log"
	"time"
	"context"
)

type Aggregator struct {
	db     *sql.DB
	config *Config
}

func fetchNews(interval int, fetchFunc func(), ctx context.Context) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fetchFunc()
		case <-ctx.Done():
			return
		}
	}
}

func (agr Aggregator) collectNewsAndSave(newsChan <-chan News, ctx context.Context) {
	for {
		select {
		case news := <-newsChan:
			err := agr.saveNews(&news)
			if err != nil {
				log.Println("error while saving news:", err)
			}
		case <-ctx.Done():
			return
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

func (agr Aggregator) Run(ctx context.Context) {
	newsChan := make(chan News, 200)

	go agr.collectNewsAndSave(newsChan, ctx)

	for _, rule := range agr.config.RSSRules {
		go fetchNews(rule.Interval, func() {
			parseRSS(rule, newsChan)
		}, ctx)
	}

	for _, rule := range agr.config.HTMLRules {
		go fetchNews(rule.Interval, func() {
			parseHTML(rule, newsChan)
		}, ctx)
	}
}
