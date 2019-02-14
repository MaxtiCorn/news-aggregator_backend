package main

import (
	"flag"
	"log"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dbPath := flag.String("db", "database.db", "path to database file")
	configPath := flag.String("config", "config.json", "path to config file")
	port := flag.String("port", "80", "port to run")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Println("create aggregator")
	aggregator, err := NewAggregator(*dbPath, *configPath)
	if err != nil {
		log.Println("error while creating aggregator", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	log.Println("run aggregator")
	aggregator.Run(ctx)

	log.Println("run api on port", *port)
	go runAPI(aggregator, *port)

	log.Println("work..")

	select {
	case <-ctx.Done():
		cancel()
	case <-sigs:
		cancel()
	}
}
