package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Routes ... routes func for all the routes
func Routes() {
	r := mux.NewRouter()

	r.HandleFunc("/api/books", GetBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", GetBook).Methods("GET")
	r.HandleFunc("/api/books", CreateBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", UpdateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", DeleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
