package main

import (
	"log"
	"net/http"

	shorneter "github.com/mikkoryynanen/url-shortener/internal"
)

func main() {
	shortener := shorneter.NewShortener()

	http.HandleFunc("/shorten", shortener.HandleShorten)
	http.HandleFunc("/s", shortener.HandleShortened)

	log.Fatal(http.ListenAndServe(":3000", nil))
}