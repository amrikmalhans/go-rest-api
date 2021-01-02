package main

import (
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/gorilla/mux"
)

// Routes ... routes func for all the routes
func Routes(JwtMiddleware *jwtmiddleware.JWTMiddleware) {
	r := mux.NewRouter()

	r.HandleFunc("/api/books", GetBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", GetBook).Methods("GET")
	r.Handle("/api/books", JwtMiddleware.Handler(CreateBook)).Methods("POST")
	r.HandleFunc("/api/books/{id}", UpdateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", DeleteBook).Methods("DELETE")
	r.HandleFunc("/login", Login).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}
