package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	Connect()
	Routes()

	r := mux.NewRouter()

	log.Fatal(http.ListenAndServe(":8000", r))
}
