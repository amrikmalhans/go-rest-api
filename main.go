package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	Connect()
	Init()
	Routes()

	r := mux.NewRouter()
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	log.Fatal(http.ListenAndServe(":8000", corsWrapper.Handler(r)))
}
