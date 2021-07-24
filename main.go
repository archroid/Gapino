package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var usersCollection *mongo.Collection
var relationsCollection *mongo.Collection

func main() {
	initDatabase()
	initRouter()
}

func initRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.HandleFunc("/auth/login", loginHandler)
	r.HandleFunc("/auth/register", registerHandler).Methods("POST")
	r.HandleFunc("/user/updateUser", updateUserHandler).Methods("POST")
	r.HandleFunc("/user/uploadImage", uploadImageHandler).Methods("POST")
	r.HandleFunc("/user/getUser", getUserHandler).Methods("POST")

	r.HandleFunc("/relation/follow", followHandler).Methods("POST")
	r.HandleFunc("/relation/unfollow", unfollowHandler).Methods("POST")

	staticDir := "/images/"
	http.Handle(staticDir, http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	http.Handle("/", handlers.LoggingHandler(os.Stdout, r))
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func initDatabase() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	db, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = db.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	usersCollection = db.Database("gapino").Collection("users")
	relationsCollection = db.Database("gapino").Collection("relations")
}
