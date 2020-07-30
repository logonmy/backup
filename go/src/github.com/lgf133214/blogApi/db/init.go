package db

import (
	"context"
	"github.com/lgf133214/blogApi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var (
	UserItems       *mongo.Collection
	VerifyCodeItems *mongo.Collection
	TagItems        *mongo.Collection
	PostItems       *mongo.Collection
	CategoryItems   *mongo.Collection
	CommentItems    *mongo.Collection
	GlobalItems     *mongo.Collection
	IpItems         *mongo.Collection
)

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	if client == nil {
		panic("mongo client is nil")
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		panic(err)
	}

	itemsInit(client)
	globalInit()
}

func itemsInit(client *mongo.Client) {
	PostItems = client.Database("blog").Collection("post_items")
	UserItems = client.Database("blog").Collection("user_items")
	VerifyCodeItems = client.Database("blog").Collection("verify_code_items")
	TagItems = client.Database("blog").Collection("tag_items")
	CategoryItems = client.Database("blog").Collection("category_items")
	CommentItems = client.Database("blog").Collection("comment_items")
	GlobalItems = client.Database("blog").Collection("global_item")
	IpItems = client.Database("blog").Collection("ip_items")
}

func globalInit() {
	one := GlobalItems.FindOne(context.TODO(), bson.D{})
	if err := one.Err(); err != nil {
		if err != mongo.ErrNoDocuments {
			panic(err)
		}

		i := model.Global{}
		i.PubTime = time.Now().Unix()
		i.TodayIps = []model.IpTmp{}
		_, err := GlobalItems.InsertOne(context.TODO(), i)
		if err != nil {
			panic(err)
		}
	}
}
