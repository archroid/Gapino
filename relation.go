package main

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type Follow struct {
	Follower_id  string
	Following_id string
}

func followHandler(w http.ResponseWriter, r *http.Request) {
	follower_id := r.FormValue("follower_id")
	following_id := r.FormValue("following_id")

	insertFollow := Follow{follower_id, following_id}
	_, err := relationsCollection.InsertOne(context.TODO(), insertFollow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		follow(follower_id, following_id)
		json.NewEncoder(w).Encode(map[string]string{"status": "done"})
	}

}

func unfollowHandler(w http.ResponseWriter, r *http.Request) {
	unfollower_id := r.FormValue("unfollower_id")
	unfollowing_id := r.FormValue("unfollowing_id")

	delete_filter := bson.D{{"follower_id", unfollower_id}, {"unfollowing_id", unfollowing_id}}

	_, err := usersCollection.DeleteOne(context.TODO(), delete_filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		unfollow(unfollower_id, unfollowing_id)
		json.NewEncoder(w).Encode(map[string]string{"status": "done"})
	}
}
