package main

import (
	"context"
	"fmt"
	"log"

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
