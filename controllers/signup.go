package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"restapi/helpers"
	"restapi/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// Signup ... func to create users
func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user)
	hashedPassword, _ := helpers.HashPassword(user.Password)
	user.Password = hashedPassword
	u2 := uuid.New()
	user.ID = uint64(u2.ID())
	fmt.Println(user.ID)
	userResult, err := helpers.UserCollection.Find(context.TODO(), bson.M{"username": user.Username})
	if err != nil {
		log.Fatal(err)
	}
	var userFiltered []bson.M
	if err = userResult.All(context.TODO(), &userFiltered); err != nil {
		log.Fatal(err)
	}
	emailResult, err := helpers.UserCollection.Find(context.TODO(), bson.M{"email": user.Email})
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
		insertUser, err := helpers.UserCollection.InsertOne(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted post with ID:", insertUser)
		user.Password = ""
		json.NewEncoder(w).Encode(user)
		w.WriteHeader(200)
	}
}
