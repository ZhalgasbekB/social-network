package internal

type User struct {
	ID     int    `json:"user_id" bson:"user_id"`
	Name   string `json:"name" bson:"name"`
	Avatar string `json:"avatar" bson:"avatar"`

	Posts         []int `json:"posts" bson:"posts"`
	LikedPosts    []int `json:"likedPosts" bson:"likedPosts"`
	Notifications []int `json:"notifications" bson:"notifications"`
}

type UserDTO struct {
	Name   string `json:"name" bson:"name"`
	Avatar string `json:"avatar" bson:"avatar"`
}

type Post struct {
	ID         int    `json:"post_id" bson:"post_id"`
	Content    string `json:"content" bson:"content"`
	Author     int    `json:"author" bson:"author"`
	LikesCount int    `json:"likes_count" bson:"likes_count"`
}

type PostDTO struct {
	Content string `json:"content" bson:"content"`
	Author  int    `json:"author" bson:"author"`
}

type UserIdDTO struct {
	UserId int `json:"user_id" bson:"user_id"`
}

type LikeDTO struct {
	UserID int `json:"user_id" bson:"user_id"`
}

type Notification struct {
	ID      int    `json:"notification_id"  bson:"notification_id"`
	Type    string `json:"type" bson:"type"`
	PostID  int    `json:"post_id" bson:"post_id"`
	LikedBy int    `json:"user_id" bson:"user_id"`
}
