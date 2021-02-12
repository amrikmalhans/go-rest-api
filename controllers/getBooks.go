package controllers

import ( 
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"restapi/models"
	"restapi/helpers"

	"go.mongodb.org/mongo-driver/bson"
)

// GetBooks ... get all books from the db
func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	var BookModel = models.Book
	var books []BookModel
	cursor, err := helpers.Collection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Finding all documents ERROR:", err)
	} else {
		for cursor.Next(ctx) {
			var result models.Book
			err := cursor.Decode(&result)

			// If there is a cursor.Decode error
			if err != nil {
				fmt.Println("cursor.Next() error:", err)
			} else {
				books = append(books, result)
			}
		}
	}
	json.NewEncoder(w).Encode(books)
}
