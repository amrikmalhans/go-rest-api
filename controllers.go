package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Author ...
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

//Book ...
type Book struct {
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

var books []Book

var collection *mongo.Collection

// BookCollection ... function to get the collection
func BookCollection(c *mongo.Database) {
	collection = c.Collection("books")
}

// GetBooks ... get all books from the db
func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	var books []Book
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Finding all documents ERROR:", err)
	} else {
		for cursor.Next(ctx) {
			var result Book
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

// GetBook ... get a single book by ID
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var result Book
	objID, _ := primitive.ObjectIDFromHex(params["id"])
	err := collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	} else {
		json.NewEncoder(w).Encode(result)
	}
}

// CreateBook ... func to create a new book
func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	json.NewEncoder(w).Encode(book)
	insertResult, err := collection.InsertOne(context.TODO(), book)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted post with ID:", insertResult)
}

// UpdateBook ... update book function
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ID, _ := primitive.ObjectIDFromHex(params["id"])
	var book Book
	filter := bson.M{"_id": ID}
	_ = json.NewDecoder(r.Body).Decode(&book)

	update := bson.D{
		{"$set", bson.D{
			{"isbn", book.Isbn},
			{"title", book.Title},
			{"author", bson.D{
				{"firstname", book.Author.Firstname},
				{"lastname", book.Author.Lastname},
			}},
		}},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(book)
}

// DeleteBook ... Delete func to delete a book
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ID, err := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": ID}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(deleteResult)
}
