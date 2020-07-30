package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
)

// TODO check site whether changed

func ip66Run() {
	// 没有上限
	url66Ip1 := "http://www.66ip.cn/mo.php?tqsl=2000"
	url66Ip2 := "http://www.66ip.cn/nmtq.php?getnum=2000"

	ch := make(chan string, 100)
	go util.GetByList([]string{url66Ip1, url66Ip2}, 1, re.Compile1, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "www.66ip.cn"}
	}
}
