package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
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

// Connect ... to the database
func Connect() {
	dotenv := goDotEnvVariable("DB_PASSWORD")

	godotenv.Load(".env")
	client, err := mongo.NewClient(options.Client().ApplyURI(
		"mongodb+srv://amrik:" + dotenv + "@cluster0.ye9cz.mongodb.net/restapi-go?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	// defer client.Disconnect(ctx)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected!")
	}
	db := client.Database("REST-API")
	userDB := client.Database("User")

	BookCollection(db)
	UserCollection(userDB)
	return
}

// Client ... redis client
var Client *redis.Client

//Init ... func to initialise redis db
func Init() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	Client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
