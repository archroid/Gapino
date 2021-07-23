package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Email            string
	Password         string
	Token            string
	Username         string
	Created_at       int64
	Name             string
	Avatar_url       string
	Header_url       string
	Bio              string
	Location         string
	Birthday         string
	Followings_count int64
	Followers_count  int64
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	email := r.FormValue("email")
	password := r.FormValue("password")

	filter := bson.M{"email": email, "password": password}

	err := usersCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"token": user.Token})
	}

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	var user User

	filter := bson.M{"email": email}

	err := usersCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err == nil {
		http.Error(w, "Used email", http.StatusBadRequest)
	} else {
		token := generateToken(email)
		insertUser := User{email, password, token, "", time.Now().Unix(), "Default", "", "", "", "", "", 0, 0}
		_, err := usersCollection.InsertOne(context.TODO(), insertUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			json.NewEncoder(w).Encode(map[string]string{"token": token})
		}

	}
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user *User

	token := r.FormValue("token")

	filter_find := bson.M{"token": token}
	_ = usersCollection.FindOne(context.TODO(), filter_find).Decode(&user)

	name := r.FormValue("name")
	if name == "" {
		name = user.Name
	}
	bio := r.FormValue("bio")
	if bio == "" {
		bio = user.Bio
	}
	location := r.FormValue("location")
	if location == "" {
		location = user.Location
	}
	birthday := r.FormValue("birthday")
	if birthday == "" {
		birthday = user.Birthday
	}
	followers_count, _ := strconv.ParseInt(r.FormValue("followers_count"), 10, 64)
	print(followers_count)
	if followers_count == 0 {
		followers_count = user.Followers_count
	}
	followings_count, _ := strconv.ParseInt(r.FormValue("followings_count"), 10, 64)
	if followings_count == 0 {
		followings_count = user.Followings_count
	}

	filter := bson.M{"token": token}

	update := bson.M{"$set": bson.M{
		"name":             name,
		"bio":              bio,
		"location":         location,
		"birthday":         birthday,
		"followers_count":  followers_count,
		"followings_count": followings_count,
	},
	}

	_, err := usersCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Panic(err)
	}

}

func generateToken(username string) string {
	type customClaims struct {
		Username string `json:username`
		jwt.StandardClaims
	}
	claims := customClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Issuer: "gapino",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, _ := token.SignedString([]byte("gapino"))

	return string(signedToken)
}