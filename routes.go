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
	r.HandleFunc("/api/books", CreateBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", UpdateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", DeleteBook).Methods("DELETE")
	r.HandleFunc("/signup", Signup).Methods("POST")
	r.HandleFunc("/login", Login).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}
