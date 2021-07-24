package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	delete_filter := bson.M{
		"follower_id":  unfollower_id,
		"following_id": unfollowing_id,
	}

	_, err := relationsCollection.DeleteOne(context.TODO(), delete_filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		unfollow(unfollower_id, unfollowing_id)
		json.NewEncoder(w).Encode(map[string]string{"status": "done"})
	}
}

func getFollowersHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	filter := bson.M{"following_id": id}

	findOptions := options.Find()
	findOptions.SetLimit(0)

	cur, err := relationsCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Println(err)
	}

	var results []*Follow
	for cur.Next(context.TODO()) {
		var elem Follow
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	json.NewEncoder(w).Encode(results)

}

func getFollowingsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	filter := bson.M{"follower_id": id}

	findOptions := options.Find()
	findOptions.SetLimit(0)

	cur, err := relationsCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Println(err)
	}

	var results []*Follow
	for cur.Next(context.TODO()) {
		var elem Follow
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	json.NewEncoder(w).Encode(results)
}
