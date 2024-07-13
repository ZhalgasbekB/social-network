package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func ConnectMongoDB(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		client.Disconnect(ctx)
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")
	return client, nil
}

func MongoClientMiddleware(client *mongo.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "mongoClient", client)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Increment(db *mongo.Database, sequenceName string) (int, error) {
	var result struct {
		Seq int `bson:"seq"`
	}
	countersCollection := db.Collection("counters")
	filter := bson.M{"_id": sequenceName}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	err := countersCollection.FindOneAndUpdate(context.TODO(), filter, update, options).Decode(&result)
	if err != nil {
		return 0, err
	}
	return result.Seq, nil
}

func WriteJSON(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(v)
}

func PostConvertor(dto PostDTO, id int) *Post {
	return &Post{
		ID:      id,
		Author:  dto.Author,
		Content: dto.Content,
	}
}

func UserConvertor(dto UserDTO, id int) *User {
	return &User{
		ID:     id,
		Name:   dto.Name,
		Avatar: dto.Avatar,
	}
}
