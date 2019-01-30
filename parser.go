package main

import (
	"log"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

func formatLink(link, host string) string {
	if !strings.Contains(link, "http") {
		return host + link
	} else {
		return link
	}
}

func parseHTML(rule HTMLRule, newsChan chan<- News) {
	log.Println("getting news from", rule.URL)

	doc, err := goquery.NewDocument(rule.URL)
	if err != nil {
		log.Println("error while getting news", err)
		return
	}

	news := News{}
	doc.Find(rule.ArticleSelector).Each(func(_ int, s *goquery.Selection) {
		news.Title = s.Find(rule.TitleSelector).Text()
		news.Description = s.Find(rule.DescriptionSelector).Text()
		if link, ok := s.Attr("href"); rule.LinkSelector == "self" && ok {
			news.Link = formatLink(link, rule.Host)
		} else if link, ok := s.Find(rule.LinkSelector).Attr("href"); ok {
			news.Link = formatLink(link, rule.Host)
		}
		
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