package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tweet struct {
	Tweet_id    string
	Text        string
	Creator_id  string
	Created_at  int64
	Likes_count int64
	Image_Url   string
	Likes       []string
}

func addTweetHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	creator_id := r.FormValue("id")
	tweet_id := xid.New().String()

	insertTweet := Tweet{tweet_id, text, creator_id, time.Now().Unix(), 0, "", []string{}}
	_, err := tweetsCollection.InsertOne(context.TODO(), insertTweet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "done"})
	}
}

func getTweetHandler(w http.ResponseWriter, r *http.Request) {
	tweet_id := r.FormValue("tweet_id")

	var tweet Tweet
	filter := bson.M{"tweet_id": tweet_id}

	err := tweetsCollection.FindOne(context.TODO(), filter).Decode(&tweet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		json.NewEncoder(w).Encode(tweet)
	}
}

func getUserTweetsHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.FormValue("user_id")

	filter := bson.M{"creator_id": user_id}

	findOptions := options.Find()
	findOptions.SetLimit(0)

	cur, err := tweetsCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Println(err)
	}

	var results []*Tweet
	for cur.Next(context.TODO()) {
		var elem Tweet
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

func updateTweetHandler(w http.ResponseWriter, r *http.Request) {
	tweet_id := r.FormValue("tweet_id")

	var tweet Tweet
	filter := bson.M{
		"tweet_id": tweet_id,
	}
	_ = tweetsCollection.FindOne(context.TODO(), filter).Decode(&tweet)

	text := r.FormValue("text")
	if text == "" {
		text = tweet.Text
	}

	update := bson.M{"$set": bson.M{
		"text": text,
	},
	}
	_, err := tweetsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Panic(err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "done"})
	}
}

func deleteTweetHandler(w http.ResponseWriter, r *http.Request) {
	tweet_id := r.FormValue("tweet_id")

	filter := bson.M{
		"tweet_id": tweet_id,
	}
	_, err := tweetsCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "done"})
	}
}

func tweetLikeHandler(w http.ResponseWriter, r *http.Request) {
	tweet_id := r.FormValue("tweet_id")
	user_id := r.FormValue("user_id")

	var tweet Tweet

	filter := bson.M{
		"tweet_id": tweet_id,
	}
	_ = tweetsCollection.FindOne(context.TODO(), filter).Decode(&tweet)

	update := bson.M{
		"$addToSet": bson.M{
			"likes": bson.M{
				"$each": []string{user_id},
			},
		},
		"$set": bson.M{
			"likes_count": tweet.Likes_count + 1,
		},
	}

	_, err := tweetsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Panic(err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "done"})
	}
}

func tweetUnLikeHandler(w http.ResponseWriter, r *http.Request) {
	tweet_id := r.FormValue("tweet_id")
	user_id := r.FormValue("user_id")

	var tweet Tweet

	filter := bson.M{
		"tweet_id": tweet_id,
	}

	_ = tweetsCollection.FindOne(context.TODO(), filter).Decode(&tweet)

	update := bson.M{
		"$pull": bson.M{"likes": user_id},
		"$set": bson.M{
			"likes_count": tweet.Likes_count - 1,
		},
	}

	_, err := tweetsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Panic(err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "done"})
	}

}

func uploadTweetImageHandler(w http.ResponseWriter, r *http.Request) {

	tweet_id := r.FormValue("tweet_id")

	file, _, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	os.MkdirAll("images/tweets", os.ModePerm)

	filepath := "images/tweets/" + tweet_id

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, _ = io.Copy(f, file)

	filter := bson.M{"tweet_id": tweet_id}

	update := bson.D{
		{"$set", bson.D{
			{"image_url", "http://127.0.0.1:8080/" + filepath},
		}},
	}

	_, errr := tweetsCollection.UpdateOne(context.TODO(), filter, update)
	if errr != nil {
		log.Print(err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "done"})
	}
}
