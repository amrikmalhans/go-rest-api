package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"restapi/helpers"
	"restapi/models"
)

// CreateCollection ... func to create a new collection
var CreateCollection = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var collection models.Collection
	_ = json.NewDecoder(r.Body).Decode(&collection)
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
	collection.UserID = userID
	fmt.Println(userID)
	json.NewEncoder(w).Encode(collection)
	fmt.Println(collection)
	// insertResult, err := collection.InsertOne(context.TODO(), book)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Inserted post with ID:", insertResult)
})
