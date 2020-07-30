package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
)

func atomontersoftRun() {
	url := "http://www.atomintersoft.com/products/alive-proxy/proxy-list"
	ch:=make(chan string, 100)
	go util.GetByDepth(url, []string{"www.atomintersoft.com"}, 2, 1, 0, re.Compile1, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "www.atomintersoft.com"}
	}
}

