package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Comment struct {
	Comment_id string
	Tweet_id   string
	Creator_id string
	Created_at int64
	Text       string
}

func addCommentHandler(w http.ResponseWriter, r *http.Request) {
	tweet_id := r.FormValue("tweet_id")
	creator_id := r.FormValue("creator_id")
	text := r.FormValue("text")

	comment_id := xid.New().String()
	time := time.Now().Unix()

	insertCommnet := Comment{comment_id, tweet_id, creator_id, time, text}

	_, err := commentsCollection.InsertOne(context.TODO(), insertCommnet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "done"})
	}
}

func updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	tweet_id := r.FormValue("tweet_id")
	creator_id := r.FormValue("creator_id")

	text := r.FormValue("text")

	filter := bson.M{"tweet_id": tweet_id, "creator_id": creator_id}

	update := bson.M{"$set": bson.M{
		"text": text,
	},
	}

	_, err := commentsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Panic(err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "done"})
	}

}
func deleteCommentHandler(w http.ResponseWriter, r *http.Request) {

	tweet_id := r.FormValue("tweet_id")
	creator_id := r.FormValue("creator_id")

	delete_filter := bson.M{"tweet_id": tweet_id, "creator_id": creator_id}

	_, err := commentsCollection.DeleteOne(context.TODO(), delete_filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "done"})
	}

}

func getCommentsByTweetHandler(w http.ResponseWriter, r *http.Request) {
	tweet_id := r.FormValue("tweet_id")

	filter := bson.M{"tweet_id": tweet_id}

	findOptions := options.Find()
	findOptions.SetLimit(0)

	cur, err := commentsCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Println(err)
	}

	var results []*Comment
	for cur.Next(context.TODO()) {
		var elem Comment
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
