package model

type Post struct {
	Id string `bson:"_id,omitempty"`

	Title   string `bson:"title"`
	AuthId  string `bson:"auth_id"`
	Content string `bson:"string"`

	Saying     string `bson:"saying"`
	Provenance string `bson:"provenance"`

	PubTime int64 `bson:"pub_time"`
	ModTime int64 `bson:"mod_time"`

	Heats  int64 `bson:"heats"`
	Views  int64 `bson:"views"`
	Thumbs int64 `bson:"thumbs"`

	CommentsId []string `bson:"comments_id"`
	TagsID     []string `bson:"tags_id"`
	CategoryId string   `bson:"category_id"`

	Delete bool `bson:"delete"`
}
