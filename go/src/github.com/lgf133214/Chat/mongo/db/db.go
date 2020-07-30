package db

import (
	"context"
	"github.com/lgf133214/Chat/mongo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var (
	MsgItems *mongo.Collection
)

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	MsgItems = client.Database("chat").Collection("msg_items")
}

func GetMsgs(lastTime int64) ([]model.Msg, bool) {
	ret := make([]model.Msg, 0, 10)
	cursor, err := MsgItems.Find(context.TODO(), bson.D{{"time", bson.D{{"$gt", lastTime - 1}}}})
	if err != nil {
		log.Println(err)
		return ret, false
	}
	err = cursor.All(context.TODO(), &ret)
	if err != nil {
		log.Println(err)
		return ret, false
	}
	return ret, true
}

func StoreMsg(msg model.Msg) bool {
	_, err := MsgItems.InsertOne(context.TODO(), msg)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
