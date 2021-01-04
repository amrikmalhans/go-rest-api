package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Routes ... routes func for all the routes
func Routes() {
	r := mux.NewRouter()
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})
	r.HandleFunc("/api/books", GetBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", GetBook).Methods("GET")
	r.HandleFunc("/api/books", CreateBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", UpdateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", DeleteBook).Methods("DELETE")
	r.HandleFunc("/signup", Signup).Methods("POST")
	r.HandleFunc("/login", Login).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", corsWrapper.Handler(r)))
}
