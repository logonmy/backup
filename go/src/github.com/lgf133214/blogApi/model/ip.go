package model

type IP struct {
	Id string `bson:"_id,omitempty"`
	IP string `bson:"ip"`
	TotalTimes int64 `bson:"total_times"`
	DaysNum  int64 `bson:"days_num"`
}
