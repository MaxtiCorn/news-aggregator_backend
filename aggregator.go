package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

type Aggregator struct {
	db     *sql.DB
	config *Config
}

func (agr Aggregator) parseRSS(rule RSSRule, newsChan chan<- News) {
	log.Println("getting news from", rule.URL)
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(rule.URL)
	news := News{}
	for _, item := range feed.Items {
		news.Title = item.Title
		news.Description = item.Description
		news.Link = item.Link

		newsChan <- news
	}
}

func (agr Aggregator) fetchRSSNews(rule RSSRule, newsChan chan<- News) {
	ticker := time.NewTicker(time.Duration(rule.Interval) * time.Second)
	defer ticker.Stop()

	for {
		go agr.parseRSS(rule, newsChan)
		_ = <-ticker.C
	}
}

func (agr Aggregator) parseHTML(rule HTMLRule, newsChan chan<- News) {
	log.Println("getting news from", rule.URL)

	doc, err := goquery.NewDocument(rule.URL)
	if err != nil {
		log.Println("error while getting news", err)
		return
	}

	news := News{}
	doc.Find(rule.ArticleSelector).Each(func(_ int, s *goquery.Selection) {
		if link, ok := s.Find("a").Attr("href"); ok {
			news.Link = rule.URL + link
		}
		news.Title = s.Find(rule.TitleSelector).Text()
		news.Description = s.Find(rule.DescriptionSelector).Text()
		newsChan <- news
	})
}

func (agr Aggregator) fetchHTMLNews(rule HTMLRule, newsChan chan<- News) {
	ticker := time.NewTicker(time.Duration(rule.Interval) * time.Second)
	defer ticker.Stop()

	for {
		go agr.parseHTML(rule, newsChan)
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
		log.Println("error while database initialization:", err)
		return nil, err
	}

	config, err := parseConfig(configPath)

	if err != nil {
		log.Println("error while parsing config file:", err)
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
		go agr.fetchRSSNews(rule, newsChan)
	}

	for _, rule := range agr.config.HTMLRules {
		go agr.fetchHTMLNews(rule, newsChan)
	}
}
