package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
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
	UserID uint64  `json:"ID"`
}

//User ...
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
	ID       uint64 `json:"ID"`
}

// TokenDetails ...
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

//AccessDetails ...
type AccessDetails struct {
	AccessUUID string
	UserID     uint64
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

// CreateToken ... func to create jwt token
func CreateToken(userid uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.New().String()

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

// CreateAuth ... func to store tokens in redis
func CreateAuth(userid uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := Client.Set(td.AccessUUID, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := Client.Set(td.RefreshUUID, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

//ExtractToken ... extract token duh!
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// VerifyToken ... to verify the token
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

//TokenValid ... to check if token is valid
func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

//ExtractTokenMetadata ... extract data to look into redis
func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}

		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     uint64(userID),
		}, nil
	}
	return nil, err
}

//FetchAuth ... fetch data from redis
func FetchAuth(authD *AccessDetails) (uint64, error) {
	userid, err := Client.Get(authD.AccessUUID).Result()
	if err != nil {
		return 0, err
	}
	fmt.Println(userid)
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return uint64(userID), nil
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
	tokenAuth, err := ExtractTokenMetadata(r)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(tokenAuth)
	userID, err := FetchAuth(tokenAuth)
	if err != nil {
		log.Fatal(err)
		return
	}
	book.UserID = userID
	fmt.Println(userID)
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
	fmt.Println(user)
	hashedPassword, _ := HashPassword(user.Password)
	user.Password = hashedPassword
	u2 := uuid.New()
	user.ID = uint64(u2.ID())
	fmt.Println(user.ID)
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
		user.Password = ""
		json.NewEncoder(w).Encode(user)
		w.WriteHeader(200)
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
		ts, err := CreateToken(result.ID)
		if err != nil {
			log.Fatal(err)
			return
		}
		saveErr := CreateAuth(result.ID, ts)
		if saveErr != nil {
			log.Fatal(err)
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		json.NewEncoder(w).Encode(tokens)

	}
}
