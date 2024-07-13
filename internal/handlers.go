package internal

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"regexp"
	"strconv"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	client, ok := r.Context().Value("mongoClient").(*mongo.Client)
	if !ok {
		http.Error(w, "MongoDB client not available", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:

		var dto UserIdDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		collection := client.Database("social-network").Collection("user")

		var user User
		idB := bson.M{"user_id": dto.UserId}

		if err := collection.FindOne(context.TODO(), idB).Decode(&user); err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, "Error retrieving user", http.StatusInternalServerError)
			}
			return
		}

		WriteJSON(w, user)

	case http.MethodPost:
		var dtoUser UserDTO
		if err := json.NewDecoder(r.Body).Decode(&dtoUser); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		userId, err := Increment(client.Database("social-network"), "user_id")
		if err != nil {
			http.Error(w, "Error generating user ID", http.StatusInternalServerError)
			return
		}

		user := UserConvertor(dtoUser, userId)

		collection := client.Database("social-network").Collection("user")
		_, err = collection.InsertOne(context.TODO(), user)
		if err != nil {
			http.Error(w, "Error saving user to database", http.StatusInternalServerError)
			return
		}

		WriteJSON(w, user)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

func Posts(w http.ResponseWriter, r *http.Request) {
	client, ok := r.Context().Value("mongoClient").(*mongo.Client)
	if !ok {
		http.Error(w, "MongoDB client not available", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:

		var user UserIdDTO
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		collection := client.Database("social-network").Collection("user")
		filter := bson.M{"user_id": user.UserId}

		var user1 User
		if err := collection.FindOne(context.Background(), filter).Decode(&user1); err != nil {
			http.Error(w, "Error retrieving user", http.StatusInternalServerError)
			return
		}

		var posts []Post
		for _, v := range user1.Posts {
			var post Post
			collection1 := client.Database("social-network").Collection("posts")
			filter1 := bson.M{"post_id": v}
			if err := collection1.FindOne(context.Background(), filter1).Decode(&post); err != nil {
				http.Error(w, "Error retrieving post", http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}

		WriteJSON(w, posts)

	case http.MethodPost:
		var dtoPost PostDTO

		if err := json.NewDecoder(r.Body).Decode(&dtoPost); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		postId, err := Increment(client.Database("social-network"), "post_id")
		if err != nil {
			http.Error(w, "Error generating user ID", http.StatusInternalServerError)
			return
		}

		post := PostConvertor(dtoPost, postId)

		collection := client.Database("social-network").Collection("posts")
		if _, err := collection.InsertOne(context.Background(), post); err != nil {
			http.Error(w, "Error saving post to database", http.StatusInternalServerError)
			return
		}

		var user User
		collection1 := client.Database("social-network").Collection("user")
		filter := bson.M{"user_id": post.Author}
		if err := collection1.FindOne(context.Background(), filter).Decode(&user); err != nil {
			http.Error(w, "Error retrieving user", http.StatusInternalServerError)
			return
		}

		user.Posts = append(user.Posts, postId)

		if _, err := collection1.UpdateOne(context.Background(), filter, bson.M{"$set": bson.M{"posts": user.Posts}}); err != nil {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func LikedUserPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		return
	}

	client, ok := r.Context().Value("mongoClient").(*mongo.Client)
	if !ok {
		http.Error(w, "MongoDB client not available", http.StatusInternalServerError)
		return
	}

	var dto UserIdDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	collection := client.Database("social-network").Collection("user")
	filter := bson.M{"user_id": dto.UserId}

	var user User
	if err := collection.FindOne(context.Background(), filter).Decode(&user); err != nil {
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	var posts []Post
	for _, v := range user.LikedPosts {
		var post Post
		collection1 := client.Database("social-network").Collection("posts")
		filter1 := bson.M{"post_id": v}
		if err := collection1.FindOne(context.Background(), filter1).Decode(&post); err != nil {
			http.Error(w, "Error retrieving post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	WriteJSON(w, posts)
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		return
	}

	client, ok := r.Context().Value("mongoClient").(*mongo.Client)
	if !ok {
		http.Error(w, "MongoDB client not available", http.StatusInternalServerError)
		return
	}

	re := regexp.MustCompile(`^/posts/(\d+)/like$`)
	matches := re.FindStringSubmatch(r.URL.Path)

	if matches == nil {
		http.NotFound(w, r)
		return
	}

	post_id, err := strconv.Atoi(matches[1])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var like LikeDTO
	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	var user User
	collection := client.Database("social-network").Collection("user")
	filter := bson.M{"user_id": like.UserID}
	if err := collection.FindOne(context.Background(), filter).Decode(&user); err != nil {

		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	user.LikedPosts = append(user.LikedPosts, post_id)

	if _, err := collection.UpdateOne(context.Background(), filter, bson.M{"$set": bson.M{"likedPosts": user.LikedPosts}}); err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	///TODO COUNT OF LIKE SHOULD BE INCREMENTED

	collection1 := client.Database("social-network").Collection("posts")
	filter1 := bson.M{"post_id": post_id}
	update := bson.M{"$inc": bson.M{"likes_count": 1}}
	if _, err := collection1.UpdateOne(context.Background(), filter1, update); err != nil {
		http.Error(w, "Error updating post", http.StatusInternalServerError)
		return
	}

	///TODO ADD MESSAGE TO NOTIFICATION

	noth_id, err := Increment(client.Database("social-network"), "notification_id")
	if err != nil {
		http.Error(w, "Error generating notification ID", http.StatusInternalServerError)
		return
	}

	notification := &Notification{
		ID:      noth_id,
		Type:    "like",
		PostID:  post_id,
		LikedBy: like.UserID,
	}

	collection2 := client.Database("social-network").Collection("notifications")
	if _, err := collection2.InsertOne(context.Background(), notification); err != nil {
		http.Error(w, "Error saving notification to database", http.StatusInternalServerError)
		return
	}

	////TODO ADD USER ID OF NOTIFICATION TO USER NOTIFICATION

	collection3 := client.Database("social-network").Collection("user")
	filter3 := bson.M{"posts": post_id}
	var user0 User

	if err := collection3.FindOne(context.Background(), filter3).Decode(&user0); err != nil {
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	user0.Notifications = append(user0.Notifications, noth_id)
	if _, err := collection3.UpdateOne(context.Background(), filter3, bson.M{"$set": bson.M{"notifications": user0.Notifications}}); err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, user0)
}

func Notifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		return
	}

	client, ok := r.Context().Value("mongoClient").(*mongo.Client)
	if !ok {
		http.Error(w, "MongoDB client not available", http.StatusInternalServerError)
		return
	}

	var user UserIdDTO
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	var userF User
	var notifications []Notification

	collection := client.Database("social-network").Collection("user")
	filter := bson.M{"user_id": user.UserId}

	if err := collection.FindOne(context.Background(), filter).Decode(&userF); err != nil {
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	for _, v := range userF.Notifications {
		var notification Notification
		collection1 := client.Database("social-network").Collection("notifications")
		filter1 := bson.M{"notification_id": v}
		if err := collection1.FindOne(context.Background(), filter1).Decode(&notification); err != nil {
			http.Error(w, "Error retrieving notification", http.StatusInternalServerError)
			return
		}
		notifications = append(notifications, notification)
	}

	WriteJSON(w, notifications)
}
