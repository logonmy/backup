package db

import (
	"context"
	"fmt"
	"github.com/lgf133214/Gin/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var (
	PostItems       *mongo.Collection
	UserItems       *mongo.Collection
	VerifyCodeItems *mongo.Collection
	TagItems        *mongo.Collection
	CategoryItems   *mongo.Collection
	CommentItems    *mongo.Collection
	TodayViews      *mongo.Collection
	TotalInfo       *mongo.Collection
)

// 数据库相关
func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	PostItems = client.Database("blog").Collection("post_items")
	UserItems = client.Database("blog").Collection("user_items")
	VerifyCodeItems = client.Database("blog").Collection("code_items")
	TagItems = client.Database("blog").Collection("tag_items")
	CategoryItems = client.Database("blog").Collection("category_items")
	CommentItems = client.Database("blog").Collection("comment_items")
	TodayViews = client.Database("blog").Collection("today_views")
	TotalInfo = client.Database("blog").Collection("total_info")
}

func GetTags() []model.Tag {
	filter := bson.D{}
	cursor, err := TagItems.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil
	}
	ret := make([]model.Tag, 10)
	err = cursor.All(context.TODO(), &ret)
	if err != nil {
		log.Println(err)
		return nil
	}
	return ret
}

func GetCategories() []model.Category {
	filter := bson.D{}
	cursor, err := CategoryItems.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil
	}
	ret := make([]model.Category, 10)
	err = cursor.All(context.TODO(), &ret)
	if err != nil {
		log.Println(err)
		return nil
	}
	return ret
}

func GetPosts(filter bson.D, opts ...*options.FindOptions) []model.Post {
	cursor, err := PostItems.Find(context.TODO(), filter, opts...)
	if err != nil {
		log.Println(err)
		return nil
	}
	ret := make([]model.Post, 10)
	err = cursor.All(context.TODO(), &ret)
	if err != nil {
		log.Println(err)
		return nil
	}
	return ret
}

func GetVerifyCode(code string, register bool) (bool, model.VerifyCode) {
	s := VerifyCodeItems.FindOne(context.TODO(), bson.D{{"register", register}, {"code", code}})
	ret := new(model.VerifyCode)
	err := s.Decode(ret)
	if err != nil {
		log.Println(err)
		return false, *ret
	}
	return true, *ret
}

func DelVerifyCode(id string) {
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
	}
	_, err = VerifyCodeItems.DeleteOne(context.TODO(), bson.D{{"_id", obj}})
	if err != nil {
		log.Println(err)
	}
}

func DelUser(id string) {
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
	}
	_, err = UserItems.DeleteOne(context.TODO(), bson.D{{"_id", obj}})
	if err != nil {
		log.Println(err)
	}
}


func GetUser(filter bson.D, opts ...*options.FindOneOptions) (bool, model.User) {
	result := UserItems.FindOne(context.TODO(), filter, opts...)
	ret := new(model.User)
	err := result.Decode(ret)
	if err != nil {
		log.Println(err)
		return false, *ret
	}
	return true, *ret
}

func GetDailyView() (bool, model.TodayViews) {
	result := TodayViews.FindOne(context.TODO(), bson.D{})
	ret := new(model.TodayViews)
	err := result.Decode(ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ret.TodayTime = time.Now().Unix()
			ret.Address=[]string{}
			_, err := TodayViews.InsertOne(context.TODO(), ret)
			if err != nil {
				log.Println(err)
				return false, *ret
			}
			return true, *ret
		}
		log.Println(err)
		return false, *ret
	}
	return true, *ret
}

func AddPost(post model.Post) bool {
	_, err := PostItems.InsertOne(context.TODO(), post)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func AddTag(tag string) bool {
	_, err := TagItems.InsertOne(context.TODO(), model.Tag{Tag: tag})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func AddCategory(category string) bool {
	_, err := CategoryItems.InsertOne(context.TODO(), model.Category{Category: category})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func AddComment(userId, postId string, comment string) bool {
	c := model.Comment{UserId: userId, PostId: postId, Content: comment}
	id, err := CommentItems.InsertOne(context.TODO(), c)
	if err != nil {
		log.Println(err)
		return false
	}
	c.Id = fmt.Sprintf("%s", id.InsertedID)

	log.Println("c.Id = " + c.Id)

	_, err = PostItems.UpdateOne(context.TODO(), bson.D{{"_id", postId}}, bson.D{{
		"$push", bson.D{{
			"comments", c,
		}},
	}})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func AddUser(user model.User) bool {
	_, err := UserItems.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func AddVerifyCode(code model.VerifyCode) bool {
	_, err := VerifyCodeItems.InsertOne(context.TODO(), code)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func DailyViewsPushAddr(addr string) {
	_, err := TodayViews.UpdateOne(context.TODO(), bson.D{}, bson.D{{
		"$push", bson.D{{
			"address", addr,
		}},
	}})
	if err != nil {
		log.Println(err)
	}
}

func DelSession(session string) bool {
	_, err := UserItems.UpdateOne(context.TODO(), bson.D{{
		"session", session,
	}}, bson.D{{
		"$set", bson.D{
			{"session", nil},
			{"expire_time", nil},
		},
	}})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func UpdateUser(id string, update bson.D) bool {
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = UserItems.UpdateOne(context.TODO(), bson.D{{"_id", obj}}, update)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
