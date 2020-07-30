package model

type IpTmp struct {
	IP         string `bson:"ip"`
	Ban        bool   `bson:"ban"`
	VisitNum   int32  `bson:"visit_num"`
	StartTime  int64  `bson:"start_time"`
	TodayTimes int32  `bson:"today_times"`
}

type Global struct {
	TotalTimes int64   `bson:"total_times"`
	PubTime    int64   `bson:"pub_time"`
	TodayIps   []IpTmp `bson:"today_ips"`
}
