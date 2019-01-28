package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var news []News

func getNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func runAPI() {
	news = append(news, News{Title: "test news 0"})
	news = append(news, News{Title: "test news 1"})

	router := mux.NewRouter()
	router.HandleFunc("/getNews", getNews).Methods("GET")
	log.Fatal(http.ListenAndServe(":69", router))
}
