package main

import (
	"QuotesAPI/internal/handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	QuoteHandler := handler.NewQuoteHandler()

	r.HandleFunc("/quotes", QuoteHandler.QuotesPOST).Methods("POST")
	r.HandleFunc("/quotes", QuoteHandler.QuotesGET).Methods("GET")
	r.HandleFunc("/quotes/random", QuoteHandler.QuotesRandom).Methods("GET")
	r.HandleFunc("/quotes/{id}", QuoteHandler.QuotesDelete).Methods("DELETE")

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
