package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"restapi/helpers"
	"restapi/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetBook ... get a single book by ID
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var result models.Book
	objID, _ := primitive.ObjectIDFromHex(params["id"])
	err := helpers.Collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	} else {
		json.NewEncoder(w).Encode(result)
	}
}
