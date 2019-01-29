package main

import "log"

func main() {
	log.Println("creating aggregator")
	aggregator, err := NewAggregator("database.db", "config.json")
	if err != nil {
		log.Println("error while creating aggregator", err)
		return
	}

	log.Println("running aggregator")
	aggregator.Run()

	log.Println("running api")
	runAPI(aggregator)
}
