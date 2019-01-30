package main

import (
	"flag"
	"log"
)

func main() {
	dbPath := flag.String("db", "database.db", "path to database file")
	configPath := flag.String("config", "config.json", "path to config file")

	log.Println("creating aggregator")
	aggregator, err := NewAggregator(*dbPath, *configPath)
	if err != nil {
		log.Println("error while creating aggregator", err)
		return
	}

	log.Println("running aggregator")
	aggregator.Run()

	log.Println("running api")
	runAPI(aggregator)
}
