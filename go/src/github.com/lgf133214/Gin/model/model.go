package model

// 模型定义
// id 的omitempty必须要有，不然分配的id都是空串

// todo 评论回复
// todo 评论删除

type Post struct {
	Id       string    `bson:"_id,omitempty"`
	UserId   string    `bson:"user_id"`
	Title    string    `bson:"title"`
	Content  string    `bson:"content"`
	Images   []string  `bson:"images"`
	PubTime  int64     `bson:"pub_time"`
	ModTime  int64     `bson:"mod_time"`
	Tags     []Tag     `bson:"tags"`
	Category Category  `bson:"category"`
	Comments []Comment `bson:"comments"`
	Views    int64     `bson:"views"`
}

type Tag struct {
	Id  string `bson:"_id,omitempty"`
	Tag string `bson:"tag"`
}

// todo 父分类，像链表一样，看心情
type Category struct {
	Id       string `bson:"_id,omitempty"`
	Category string `bson:"category"`
}

type User struct {
	Id            string `bson:"_id,omitempty"`
	Name          string `bson:"name"`
	PassWord      string `bson:"password"`
	Validated     bool   `bson:"validated"`
	SuperUser     bool   `bson:"super_user"`
	Mail          string `bson:"mail"`
	JoinTime      int64  `bson:"join_time"`
	Icon          string `bson:"icon"`
	LastLoginTime int64  `bson:"last_login_time"`
	Session       string `bson:"session"`
	ExpireTime    int64  `bson:"expire_time"`
}

type VerifyCode struct {
	Id       string `bson:"_id,omitempty"`
	UserId   string `bson:"user_id"`
	Register bool   `bson:"register"`

	Code       string `bson:"code"`
	ExpireTime int64  `bson:"expire_time"`
}

type Comment struct {
	Id        string   `bson:"_id,omitempty"`
	UserId    string   `bson:"user_id"`
	PostId    string   `bson:"post_id"`
	CommentId string   `bson:"comment_id"`
	Content   string   `bson:"content"`
	Time      int64    `bson:"time"`
	Thumbs    []string `bson:"thumbs"`
}

type TodayViews struct {
	Id        string   `bson:"_id,omitempty"`
	TodayTime int64    `bson:"today_time"`
	Address   []string `bson:"address"`
}

type TotalInfo struct {
	Id        string `bson:"_id,omitempty"`
	Views     int64  `bson:"views"`
	StartTime int64  `bson:"start_time"`
}
