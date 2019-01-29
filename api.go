package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (agr Aggregator) getNewsHandler(w http.ResponseWriter, r *http.Request) {
	news, err := agr.getAllNews()
	if err != nil {
		log.Println("error while getting all news", err)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func (agr Aggregator) searchNewsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	news, err := agr.searchNews(vars["search"])
	if err != nil {
		log.Println("error while searching news", err)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func runAPI(agr *Aggregator) {
	router := mux.NewRouter()
	router.HandleFunc("/getNews", agr.getNewsHandler).Methods("GET")
	router.HandleFunc("/searchNews", agr.searchNewsHandler).Queries("search", "{search}").Methods("GET")
	log.Fatal(http.ListenAndServe(":69", router))
}
