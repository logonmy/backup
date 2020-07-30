package model

type VerifyCode struct {
	Id         string `bson:"_id,omitempty"`
	Email     string `bson:"email"`
	Register   bool   `bson:"register"`
	Content    string `bson:"content"`
	ExpireTime int64  `bson:"expire_time"`
}
