package controllers
 
import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"restapi/helpers"
	"restapi/models"

	"go.mongodb.org/mongo-driver/bson"
)

// Login ... func to log users in
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	var result models.User
	err := helpers.UserCollection.FindOne(context.TODO(), bson.D{{Key: "email", Value: user.Email}}).Decode(&result)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(`{"message": "Email Invalid :(, try something else"}`))
	} else {
		userPassword := helpers.ComparePassword(result.Password, user.Password)
		if userPassword != true {
			w.WriteHeader(401)
			w.Write([]byte(`{"message": "Wrong Password :(, try something else"}`))
			return
		}
		ts, err := helpers.CreateToken(result.ID)
		if err != nil {
			log.Fatal(err)
			return
		}
		saveErr := helpers.CreateAuth(result.ID, ts)
		if saveErr != nil {
			log.Fatal(err)
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		cookie := &http.Cookie{Name: "refresh_token", Value: ts.RefreshToken, HttpOnly: true}
		http.SetCookie(w, cookie)
		json.NewEncoder(w).Encode(tokens)

	}
}
