package db

import (
	"context"
	"github.com/lgf133214/blogApi/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func findAllPost() {
	PostItems.Find(context.TODO(), bson.D{}, options.Find().SetLimit(config.PerPageNum).SetSort(bson.D{
		{"pub_time", -1},
	}))

}