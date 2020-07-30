package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
)

func proxyIpListRun() {
	url := "http://proxy-ip-list.com/"

	ch:=make(chan string, 100)

	go util.GetByDepth(url, []string{"proxy-ip-list.com"}, 2, 1, 0, re.Compile1, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "proxy-ip-list.com"}
	}
}
