package model

type Tag struct {
	Id      string   `bson:"_id,omitempty"`
	Content string   `bson:"content"`
	PostsId []string `bson:"posts_id"`
}
