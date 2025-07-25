package main

import (
	"log"
	"net/http"

	"github.com/grayankit/go_url_shortener/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Home)
	mux.HandleFunc("/shorten", handler.Shorten)
	mux.HandleFunc("/u/", handler.Redirect)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
