package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Id               string
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
		http.Error(w, "This email is already available in the database", http.StatusBadRequest)
	} else {
		token := generateToken(email)
		id := xid.New().String()
		createdAt := time.Now().Unix()
		insertUser := User{id, email, password, token, "", createdAt, "Default", "", "", "", "", "", 0, 0}
		_, err := usersCollection.InsertOne(context.TODO(), insertUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			json.NewEncoder(w).Encode(map[string]string{"token": token, "id": id})
		}

	}
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user *User

	token := r.FormValue("token")

	filter := bson.M{"token": token}
	_ = usersCollection.FindOne(context.TODO(), filter).Decode(&user)

	email := r.FormValue("email")
	if email == "" {
		email = user.Email
	}

	password := r.FormValue("password")
	if password == "" {
		password = user.Password
	}

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

	update := bson.M{"$set": bson.M{
		"email":    email,
		"password": password,
		"name":     name,
		"bio":      bio,
		"location": location,
		"birthday": birthday,
	},
	}

	_, err := usersCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Panic(err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "done"})
	}
}

func uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	id := r.FormValue("id")
	imgType := r.FormValue("type")

	file, _, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var filepath string

	os.MkdirAll("images/avatars", os.ModePerm)
	os.MkdirAll("images/headers", os.ModePerm)

	switch imgType {
	case "avatar":
		filepath = "images/avatars/" + id
		update := bson.D{
			{"$set", bson.D{
				{"avatar_url", "http://192.168.1.108:8080/" + filepath},
			}},
		}
		filter := bson.M{"token": token}
		_, err := usersCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Print(err)
		}
	case "header":
		filepath = "images/headers/" + id
		update := bson.D{
			{"$set", bson.D{
				{"header_url", "http://127.0.0.1:8080/" + filepath},
			}},
		}
		filter := bson.M{"token": token}
		_, err := usersCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Print(err)
		}
	}

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, _ = io.Copy(f, file)

	json.NewEncoder(w).Encode(map[string]string{"status": "http://127.0.0.1:8080/" + filepath})
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	id := r.FormValue("id")
	filter_find := bson.M{"id": id}
	err := usersCollection.FindOne(context.TODO(), filter_find).Decode(&user)
	if err != nil {
		http.Error(w, "Could not find the user", http.StatusBadRequest)
	} else {
		json.NewEncoder(w).Encode(user)
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

func follow(follower_id string, following_id string) {
	var user User
	filter_follower := bson.M{"id": follower_id}
	_ = usersCollection.FindOne(context.TODO(), filter_follower).Decode(&user)
	update_followers := bson.M{"$set": bson.M{
		"followings_count": user.Followings_count + 1,
	},
	}
	_, _ = usersCollection.UpdateOne(context.TODO(), filter_follower, update_followers)

	filter_following := bson.M{"id": following_id}
	_ = usersCollection.FindOne(context.TODO(), following_id).Decode(&user)
	update_followings := bson.M{"$set": bson.M{
		"followers_count": user.Followers_count + 1,
	},
	}
	_, _ = usersCollection.UpdateOne(context.TODO(), filter_following, update_followings)
}

func unfollow(unfollower_id string, unfollowing_id string) {
	var user User

	filter_unfollower := bson.M{"id": unfollower_id}
	_ = usersCollection.FindOne(context.TODO(), filter_unfollower).Decode(&user)
	update_followers := bson.M{"$set": bson.M{
		"followings_count": user.Followings_count - 1,
	},
	}
	_, _ = usersCollection.UpdateOne(context.TODO(), filter_unfollower, update_followers)

	filter_unfollowing := bson.M{"id": unfollowing_id}
	_ = usersCollection.FindOne(context.TODO(), filter_unfollowing).Decode(&user)
	update_followings := bson.M{"$set": bson.M{
		"followers_count": user.Followers_count - 1,
	},
	}
	_, _ = usersCollection.UpdateOne(context.TODO(), filter_unfollowing, update_followings)

}
