package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func timeLineHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.FormValue("user_id")

	var tweets []Tweet
	followingsIds := getFollowingIds(user_id)

	// Get user's timeline
	filter := bson.M{"creator_id": bson.M{"$in": followingsIds}}

	findOptions := options.Find()
	findOptions.SetLimit(0)

	cur, err := tweetsCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Println("Error:", err)
	}
	for cur.Next(context.TODO()) {
		tweets = append(tweets, Tweet{})
		err = cur.Decode(&tweets[len(tweets)-1])
		if err != nil {
			log.Println("Error:", err)
		}
	}
	if err := cur.Err(); err != nil {
		log.Println("Error:", err)
	}
	cur.Close(context.TODO())

	json.NewEncoder(w).Encode(tweets)

}

func getFollowingIds(user_id string) []string {

	var followings []*Follow
	var followings_id []string

	filter := bson.M{"follower_id": user_id}
	findOptions := options.Find()
	findOptions.SetLimit(0)
	cur, err := relationsCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var elem Follow
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		followings = append(followings, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())

	// Get Followings id list
	for _, elem := range followings {
		followings_id = append(followings_id, elem.Following_id)
	}

	return followings_id
}
