package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

//User ...
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var books []Book

var collection *mongo.Collection
var userCollection *mongo.Collection

// BookCollection ... function to get the collection
func BookCollection(c *mongo.Database) {
	collection = c.Collection("books")
}

// UserCollection ... function to get the collection
func UserCollection(c *mongo.Database) {
	userCollection = c.Collection("users")
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
var CreateBook = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	json.NewEncoder(w).Encode(book)
	fmt.Println(book)
	insertResult, err := collection.InsertOne(context.TODO(), book)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted post with ID:", insertResult)
})

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

// ComparePassword hashes the test password and then compares
// the two hashes.
func ComparePassword(hashBase64, testPassword string) bool {

	// Decode the real hashed and salted password so we can
	// split out the salt
	hashBytes, err := base64.StdEncoding.DecodeString(hashBase64)
	if err != nil {
		fmt.Println("Error, we were given invalid base64 string", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword(hashBytes, []byte(testPassword))
	return err == nil
}

// HashPassword hashes the clear-text password and encodes it as base64,
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10 /*cost*/)
	if err != nil {
		return "", err
	}

	// Encode the entire thing as base64 and return
	hashBase64 := base64.StdEncoding.EncodeToString(hashedBytes)

	return hashBase64, nil
}

// Signup ... func to create users
func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	hashedPassword, _ := HashPassword(user.Password)
	user.Password = hashedPassword

	userResult, err := userCollection.Find(context.TODO(), bson.M{"username": user.Username})
	if err != nil {
		log.Fatal(err)
	}
	var userFiltered []bson.M
	if err = userResult.All(context.TODO(), &userFiltered); err != nil {
		log.Fatal(err)
	}
	emailResult, err := userCollection.Find(context.TODO(), bson.M{"email": user.Email})
	if err != nil {
		log.Fatal(err)
	}
	var emailFiltered []bson.M
	if err = emailResult.All(context.TODO(), &emailFiltered); err != nil {
		log.Fatal(err)
	}
	if len(userFiltered) > 0 {
		w.WriteHeader(409)
		w.Write([]byte(`{"message": "Username Already exists :(, try something else"}`))
	} else if len(emailFiltered) > 0 {
		w.WriteHeader(409)
		w.Write([]byte(`{"message": "Email Already exists :("}`))
	} else {
		insertUser, err := userCollection.InsertOne(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted post with ID:", insertUser)
		json.NewEncoder(w).Encode(user)
	}
}

// Login ... func to log users in
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	var result User
	err := userCollection.FindOne(context.TODO(), bson.D{{Key: "email", Value: user.Email}}).Decode(&result)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(`{"message": "Email Invalid :(, try something else"}`))
	} else {
		userPassword := ComparePassword(result.Password, user.Password)
		if userPassword != true {
			w.WriteHeader(401)
			w.Write([]byte(`{"message": "Wrong Password :(, try something else"}`))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"message": "successfully logged in!"}`))

	}

	// } else if len(emailFiltered) < 0 {

	// } else {
	// 	w.WriteHeader(200)
	// 	w.Write([]byte(`{"message": "Welcome :)"}`))
	// }
}
