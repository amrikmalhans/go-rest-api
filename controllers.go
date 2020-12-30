package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Author ...
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

//Book ...
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

var collection *mongo.Collection

// BookCollection ... function to get the collection
func BookCollection(c *mongo.Database) {
	collection = c.Collection("books")
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookss := []Book{}
	err := collection.FindOne(context.TODO(), bson.M{}).Decode(&bookss)

	if err != nil {

		log.Fatal(err)

	}
	fmt.Println(bookss)

	json.NewEncoder(w).Encode(bookss)

	return
}

// InsertBook ... function to insert book to the db
func InsertBook(title string, isbn string, id string, fname string, lname string) {

	books := Book{
		title,
		isbn,
		id,
		&Author{fname, lname},
	}

	insertResult, err := collection.InsertOne(context.TODO(), books)

	if err != nil {

		log.Fatal(err)

	}

	fmt.Println("Inserted post with ID:", insertResult.InsertedID)
}
