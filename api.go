package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (agr Aggregator) getNewsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	news, err := agr.getNews(vars["count"], vars["offset"])
	if err != nil {
		log.Println("error while getting news", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func (agr Aggregator) getNewsBySourceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	news, err := agr.getNewsBySource(vars["source"])
	if err != nil {
		log.Println("error while searching news", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func runAPI(agr *Aggregator, port string) {
	router := mux.NewRouter()
	router.HandleFunc("/getNews", agr.getNewsHandler).Queries("count", "{count}", "offset", "{offset}").Methods("GET")
	router.HandleFunc("/searchNews", agr.searchNewsHandler).Queries("search", "{search}").Methods("GET")
	router.HandleFunc("/getNewsBySource", agr.getNewsBySourceHandler).Queries("source", "{source}").Methods("GET")
	log.Fatal(http.ListenAndServe(":" + port, router))
}
