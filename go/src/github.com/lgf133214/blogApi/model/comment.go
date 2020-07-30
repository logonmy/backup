package model

type Comment struct {
	Id string `bson:"_id,omitempty"`

	PostId         string `bson:"post_id"`
	UserId         string `bson:"user_id"`
	ReplyCommentId string `bson:"reply_comment_id"`
	ReplyUserId    string `bson:"reply_user_id"`

	Content string `bson:"content"`
	PubTime int64  `bson:"pub_time"`
	Thumbs  int64  `bson:"thumbs"`

	Delete bool `bson:"delete"`
}
