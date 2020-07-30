package storage

import (
	"context"

	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var (
	successItems *mongo.Collection
	bufferItems  *mongo.Collection
	unableItems  *mongo.Collection
	// store the new come in proxy ip temporally. spider -> mongo
	BufferChan = make(chan model.BufferItem, 1000)
	// validate the buffer proxy ip. mongo -> validator
	BufferValChan = make(chan model.BufferItem, 1000)
	// validate the available proxy ip. mongo -> validator
	SuccessValChan = make(chan model.SuccessItem, 1000)
	// store unable proxy from success item temporary, try time>5, remove
	UnableValChan = make(chan model.UnableItem, 1000)
)

func init() {
	// connect to mongo
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println("mongo init err")
		return
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Println("mongo init err")
		return
	}

	// get the collection connect, it will alive until exit
	successItems = client.Database("proxies").Collection("success_items")
	bufferItems = client.Database("proxies").Collection("buffer_items")
	unableItems = client.Database("proxies").Collection("unable_items")
}

// store a proxy to buffer collection, if exist or occur error, return false
func StoreToBufferDB(item model.BufferItem) {
	ok := findFromBufferDB(item.ProxyIp)
	if ok {
		return
	}
	_, err := bufferItems.InsertOne(context.TODO(), item)
	if err != nil {
		log.Println("store error" + err.Error())
		return
	}
	logger.Logger.Debug("insert " + item.ProxyIp + " " + item.Source + " to buffer collection success")
	return
}

// check the proxy whether exist in buffer collection
func findFromBufferDB(proxyIp string) bool {
	filter := bson.D{{"proxy_ip", proxyIp}}
	err := bufferItems.FindOne(context.TODO(), filter).Err()
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Println("unexpected err " + err.Error())
		}
		return false
	}

	return true
}

// item must be init, otherwise will be the default value
func StoreToSuccessDB(item model.SuccessItem) {
	ok := findFromSuccessDB(item.ProxyIp)
	if ok {
		filter := bson.D{{"proxy_ip", item.ProxyIp}}

		update := bson.D{{"$set",
			bson.D{
				{"verify_time", item.VerifyTime},
				{"http", item.Http},
				{"google", item.Google},
				{"socks5", item.Socks5},
				{"https", item.Https},
			}},
		}
		_, err := successItems.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Println(err)
		}
		return
	}
	_, err := successItems.InsertOne(context.TODO(), item)
	if err != nil {
		log.Println(err)
	}
	logger.Logger.Debug("validate " + item.ProxyIp + " " + item.Source + " success")
}

// check the proxy whether exist in success collection
func findFromSuccessDB(proxyIp string) bool {
	filter := bson.D{{"proxy_ip", proxyIp}}
	err := successItems.FindOne(context.TODO(), filter).Err()
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Println(err)
		}
		return false
	}
	return true
}

func StoreToUnableDB(item model.UnableItem) {
	i, ok := findFromUnableDB(item.ProxyIp)
	if ok {
		if i.TryTimes >= 5 {
			RemoveUnableItem(item.ProxyIp)
			return
		}
		filter := bson.D{{"proxy_ip", item.ProxyIp}}

		update := bson.D{{"$inc", bson.D{{"try_times", 1}}}}
		_, err := unableItems.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Println(err)
		}
		return
	}
	_, err := unableItems.InsertOne(context.TODO(), item)
	if err != nil {
		log.Println(err)
	}
	logger.Logger.Debug("insert " + item.ProxyIp + " " + item.Source + " to unable collection success")
}

func findFromUnableDB(proxyIp string) (model.UnableItem, bool) {
	var ret = model.UnableItem{}
	filter := bson.D{{"proxy_ip", proxyIp}}
	res := unableItems.FindOne(context.TODO(), filter)

	err := res.Err()
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Println(err)
		}
		return ret, false
	}
	err = res.Decode(&ret)
	if err != nil {
		log.Println(err)
		return ret, false
	}

	return ret, true
}

// send value to bufferValChan when the chan is empty
func GetBufferItem() bool {
	select {
	case <-BufferValChan:
		return false
	default:
	}
	count, err := bufferItems.CountDocuments(context.TODO(), bson.D{})

	if err != nil {
		log.Println(err)
		return false
	}
	if count <= 0 {
		return true
	}

	cursor, err := bufferItems.Find(context.TODO(), bson.D{}, options.Find().SetLimit(1000))
	if err != nil {
		log.Println(err)
		return false
	}

	for cursor.Next(context.TODO()) {
		var elem model.BufferItem
		err := cursor.Decode(&elem)
		if err != nil {
			log.Println(err)
			continue
		}
		BufferValChan <- elem
	}

	// Close the cursor once finished
	cursor.Close(context.TODO())
	return false
}

// send success item to success chan sort by validate time
func GetSuccessItem() bool {
	select {
	case <-SuccessValChan:
		return false
	default:
	}
	filter := bson.D{{"verify_time", bson.D{{"$lt", time.Now().Unix() - 250}}}}

	count, err := successItems.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return false
	}
	if count <= 0 {
		return true
	}

	cursor, err := successItems.Find(context.TODO(), filter, options.Find().SetLimit(1000))
	if err != nil {
		log.Println(err)
		return false
	}

	for cursor.Next(context.TODO()) {
		var elem model.SuccessItem
		err := cursor.Decode(&elem)
		if err != nil {
			log.Println(err)
			continue
		}
		SuccessValChan <- elem
	}

	// Close the cursor once finished
	cursor.Close(context.TODO())
	return false
}

func GetUnableItem() bool {
	select {
	case <-UnableValChan:
		return false
	default:
	}

	count, err := unableItems.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		log.Println(err)
		return false
	}
	if count <= 0 {
		return true
	}

	cursor, err := unableItems.Find(context.TODO(), bson.D{}, options.Find().SetLimit(1000))
	if err != nil {
		log.Println(err)
		return false
	}

	for cursor.Next(context.TODO()) {
		var elem model.UnableItem
		err := cursor.Decode(&elem)
		if err != nil {
			log.Println(err)
			continue
		}
		UnableValChan <- elem
	}

	// Close the cursor once finished
	cursor.Close(context.TODO())
	return false
}

func RemoveBufferItem(proxyIp string) {
	_, err := bufferItems.DeleteOne(context.TODO(), bson.D{{"proxy_ip", proxyIp}})
	if err != nil {
		log.Println(err)
	}
}

func RemoveSuccessItem(proxyIp string) {
	_, err := successItems.DeleteOne(context.TODO(), bson.D{{"proxy_ip", proxyIp}})
	if err != nil {
		log.Println(err)
	}
}

func RemoveUnableItem(proxyIp string) {
	_, err := unableItems.DeleteOne(context.TODO(), bson.D{{"proxy_ip", proxyIp}})
	if err != nil {
		log.Println(err)
	}
}

func GetSuccessItems(skip int64, paras map[string][]string) (int64, []model.SuccessItem, bool) {
	filter := bson.M{}

	if util.GetPara(paras, "http") != "" && (util.GetPara(paras, "http") == "1" || util.GetPara(paras, "http") == "0") {
		if util.GetPara(paras, "http") == "1" {
			filter["http"] = true
		} else {
			filter["http"] = false
		}
	}
	if util.GetPara(paras, "https") != "" && (util.GetPara(paras, "https") == "1" || util.GetPara(paras, "https") == "0") {
		if util.GetPara(paras, "https") == "1" {
			filter["https"] = true
		} else {
			filter["https"] = false
		}
	}
	if util.GetPara(paras, "google") != "" && (util.GetPara(paras, "google") == "1" || util.GetPara(paras, "google") == "0") {
		if util.GetPara(paras, "google") == "1" {
			filter["google"] = true
		} else {
			filter["google"] = false
		}
	}
	if util.GetPara(paras, "socks5") != "" && (util.GetPara(paras, "socks5") == "1" || util.GetPara(paras, "socks5") == "0") {
		if util.GetPara(paras, "socks5") == "1" {
			filter["socks5"] = true
		} else {
			filter["socks5"] = false
		}
	}
	if util.GetPara(paras, "anonymous") != "" && (util.GetPara(paras, "anonymous") == "2" || util.GetPara(paras, "anonymous") == "1" || util.GetPara(paras, "anonymous") == "0") {
		filter["anonymous"] = util.GetPara(paras, "anonymous")
	}
	count, err := successItems.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return 0, nil, false
	}

	if count <= skip {
		return 0, nil, false
	}

	cursor, err := successItems.Find(context.TODO(), filter, options.Find().SetSkip(skip).SetLimit(20).SetSort(bson.D{{"verify_time", -1}}))
	if err != nil {
		log.Println(err)
		return 0, nil, false
	}
	defer cursor.Close(context.TODO())
	ret := make([]model.SuccessItem, 0, 20)
	err = cursor.All(context.TODO(), &ret)
	if err != nil {
		log.Println(err)
		return 0, nil, false
	}

	return count, ret, true
}

func GetSuccessItemsForRPC() ([]model.SuccessItem, bool) {
	filter := bson.D{{"verify_time", -1}}

	cursor, err := successItems.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	defer cursor.Close(context.TODO())
	ret := make([]model.SuccessItem, 0, 300)
	err = cursor.All(context.TODO(), &ret)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return ret, true
}
