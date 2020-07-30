package model

type User struct {
	Id string `bson:"_id,omitempty"`

	NickName string `bson:"nick_name"`
	Email    string `bson:"email"`
	Icon     string `bson:"icon"`
	Password string `bson:"password"`

	RegisterTime int64 `bson:"register_time"`

	Session       string `bson:"session"`
	ExpireTime    int64  `bson:"expire_time"`
	LastLoginTime int64  `bson:"last_login_time"`

	PostsId    []string `bson:"posts_id"`
	CommentsId []string `bson:"comments_id"`

	BadStatus           bool  `bson:"bad_status"`
	BadStatusExpireTime int64 `bson:"bad_status_expire_time"`

	SuperUser bool `bson:"superuser"`
	Admin     bool `bson:"admin"`

	SendTimes int32 `bson:"send_times"`
	SendTime  int64 `bson:"send_time"`

	Saying     string `bson:"saying"`
	Provenance string `bson:"provenance"`
	Profile    string `bson:"profile"`

	Subscript     bool  `bson:"subscript"`
	SubscriptTime int64 `bson:"subscript_time"`
}
