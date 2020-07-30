package model

type Msg struct {
	Id       string `bson:"_id,omitempty"`
	Content  string `bson:"content"`
	UserName string `bson:"user_name"`
	Time     int64  `bson:"time"`
}

type Global struct {
	Id           string `bson:"_id, omitempty"`
	LastSendTime int64  `bson:"last_send_time"`
}
