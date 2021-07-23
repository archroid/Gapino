package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Email    string
	Password string
	Token    string
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
		insertUser := User{email, password, token}
		_, err := usersCollection.InsertOne(context.TODO(), insertUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			json.NewEncoder(w).Encode(map[string]string{"token": token})
		}

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
