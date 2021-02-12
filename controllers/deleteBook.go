package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"restapi/helpers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteBook ... Delete func to delete a book
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ID, err := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": ID}
	deleteResult, err := helpers.Collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(deleteResult)
}
