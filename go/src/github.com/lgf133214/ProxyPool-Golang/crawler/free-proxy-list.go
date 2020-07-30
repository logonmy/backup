package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
	"strconv"
)

func freeProxyListRun() {
	ch := make(chan string, 100)
	url := "https://free-proxy-list.com/?page="

	urls := make([]string, 0, 10)
	for i := range [10][0]int{} {
		urls = append(urls, url+strconv.Itoa(i+1))
	}

	go util.GetByList(urls, 1, re.Compile1, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "free-proxy-list.com"}
	}
}
