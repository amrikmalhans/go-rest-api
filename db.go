package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func db() {
	dotenv := goDotEnvVariable("DB_PASSWORD")

	godotenv.Load(".env")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://amrik:"+dotenv+"@cluster0.ye9cz.mongodb.net/restapi-go?retryWrites=true&w=majority",
	))
	collection := client.Database("testing").Collection("numbers")
	res, err := collection.InsertOne(ctx, bson.M{"name": "mmao", "value": 3.14159})
	id := res.InsertedID
	println(id)
	if err != nil {
		log.Fatal(err)
	}

}
