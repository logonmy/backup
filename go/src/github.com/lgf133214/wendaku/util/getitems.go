package util

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type SuccessItem struct {
	// ip:port
	ProxyIp    string `bson:"proxy_ip"`
	VerifyTime int64  `bson:"verify_time"`
	Google     bool   `bson:"google"`
	Http       bool   `bson:"http"`
	Https      bool   `bson:"https"`
	Socks5     bool   `bson:"socks5"`
	Source     string `bson:"source"`
	Anonymous  int    `bson:"anonymous"`
}

var successItems *mongo.Collection

func init() {
	// connect to mongo
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		panic(err)
	}

	// get the collection connect, it will alive until exit
	successItems = client.Database("proxies").Collection("success_items")
}

func GetSuccessItems() ([]SuccessItem, bool) {
	filter := bson.D{
		{"verify_time", bson.D{{"$gt", time.Now().Unix() - 180}}},
		{"https", true}}
	count, err := successItems.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, false
	}

	if count <= 0 {
		return nil, false
	}

	cursor, err := successItems.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{"verify_time", -1}}))
	if err != nil {
		return nil, false
	}
	ret := make([]SuccessItem, 0, count)
	err = cursor.All(context.TODO(), &ret)
	if err != nil {
		return nil, false
	}

	return ret, true
}
