package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
)

func proxyListDaily() {
	url := []string{"https://www.proxylistdaily.net/"}
	ch := make(chan string, 100)

	go util.GetByList(url, 1, re.Compile1, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "www.proxylistdaily.net"}
	}
}
