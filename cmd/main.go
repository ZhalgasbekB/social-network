package main

import (
	"context"
	"log"
	"net/http"
	"social-network-test/internal"
)

func main() {
	client, err := internal.ConnectMongoDB("mongodb://root:0000@mongo:27017/") // Local: mongodb://localhost:27017 mongodb://root:example@mongo:27017/
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}
	defer client.Disconnect(context.TODO())

	mux := http.NewServeMux()
	mux.Handle("/notifications", internal.MongoClientMiddleware(client)(http.HandlerFunc(internal.Notifications)))
	mux.Handle("/posts/", internal.MongoClientMiddleware(client)(http.HandlerFunc(internal.LikePost)))            // ????
	mux.Handle("/posts/liked", internal.MongoClientMiddleware(client)(http.HandlerFunc(internal.LikedUserPosts))) // ????
	mux.Handle("/profile", internal.MongoClientMiddleware(client)(http.HandlerFunc(internal.Profile)))
	mux.Handle("/posts", internal.MongoClientMiddleware(client)(http.HandlerFunc(internal.Posts)))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err.Error())
	}
}
