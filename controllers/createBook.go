package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"restapi/helpers"
	"restapi/models"
)

// CreateBook ... func to create a new book
var CreateBook = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book models.Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	tokenAuth, err := helpers.ExtractTokenMetadata(r)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(tokenAuth)
	userID, err := helpers.FetchAuth(tokenAuth)
	if err != nil {
		log.Fatal(err)
		return
	}
	book.UserID = userID
	fmt.Println(userID)
	json.NewEncoder(w).Encode(book)
	fmt.Println(book)
	insertResult, err := helpers.Collection.InsertOne(context.TODO(), book)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted post with ID:", insertResult)
})
