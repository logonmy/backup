package model

// item has passed one of the test
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

// buffer of new come in proxies
type BufferItem struct {
	// ip:port
	ProxyIp string `bson:"proxy_ip"`
	Source  string `bson:"source"`
}

// try_times > 5 remove
type UnableItem struct {
	// ip:port
	ProxyIp  string `bson:"proxy_ip"`
	Source   string `bson:"source"`
	TryTimes int    `bson:"try_times"`
}
