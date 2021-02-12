package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"restapi/helpers"
	"restapi/models"
	"strconv" 

	jwt "github.com/dgrijalva/jwt-go"
)

// Refresh ... refresh the tokens
func Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("refresh_token")
	fmt.Println(cookie)
	w.Header().Set("Content-Type", "application/json")
	var refreshToken models.RefreshToken
	refreshToken.RefreshToken = cookie.Value
	_ = json.NewDecoder(r.Body).Decode(&refreshToken)
	fmt.Println(refreshToken.RefreshToken)
	// mapToken := map[string]string{}
	// refreshToken := mapToken["refresh_token"]

	// fmt.Println(mapToken, refreshToken)
	//verify the token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	token, err := jwt.Parse(refreshToken.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		log.Fatal(err)
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		log.Fatal(err)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			log.Fatal(err)
			return
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			log.Fatal(err)
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := helpers.DeleteAuth(refreshUUID)
		if delErr != nil || deleted == 0 { //if any goes wrong
			log.Fatal(err)
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := helpers.CreateToken(userID)
		if createErr != nil {
			log.Fatal(err)
			return
		}
		//save the tokens metadata to redis
		saveErr := helpers.CreateAuth(userID, ts)
		if saveErr != nil {
			log.Fatal(saveErr)
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		cookie := &http.Cookie{Name: "refresh_token", Value: ts.RefreshToken, HttpOnly: true}
		http.SetCookie(w, cookie)
		json.NewEncoder(w).Encode(tokens)
	} else {
		json.NewEncoder(w).Encode(http.StatusUnauthorized)
	}
}
