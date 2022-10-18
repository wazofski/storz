package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Mongo() {
	/* Connect to my cluster */
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	defer cancel()

	/* List databases */
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(databases)

	/* Define my document struct */
	type Post struct {
		Title string `bson:"title,omitempty"`
		Body  string `bson:"body,omitempty"`
	}

	/* Get my collection instance */
	collection := client.Database("blog").Collection("posts")

	/* Insert documents */
	docs := []interface{}{
		bson.D{{Key: "title", Value: "World"}, {Key: "body", Value: "Hello World"}},
		bson.D{{Key: "title", Value: "Mars"}, {Key: "body", Value: "Hello Mars"}},
		bson.D{{Key: "title", Value: "Pluto"}, {Key: "body", Value: "Hello Pluto"}},
	}

	res, insertErr := collection.InsertMany(ctx, docs)
	if insertErr != nil {
		log.Fatal(insertErr)
	}

	fmt.Println(res)

	/* Iterate a cursor and print it */
	cur, currErr := collection.Find(ctx, bson.D{})

	if currErr != nil {
		panic(currErr)
	}

	defer cur.Close(ctx)

	var posts []Post
	if err = cur.All(ctx, &posts); err != nil {
		panic(err)
	}

	fmt.Println(posts)
}
