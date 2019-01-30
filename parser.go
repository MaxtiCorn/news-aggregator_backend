package main

import (
	"log"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

func parseHTML(rule HTMLRule, newsChan chan<- News) {
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

func parseRSS(rule RSSRule, newsChan chan<- News) {
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