package model

type Category struct {
	Id      string   `bson:"_id,omitempty"`
	Content string   `bson:"content"`
	PostsId []string `bson:"post_id"`
}
