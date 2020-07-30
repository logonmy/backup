package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
)

func cnProxyRun() {
	urlCnProxy := []string{"http://cn-proxy.com/", "http://cn-proxy.com/archives/218"}
	ch := make(chan string, 100)
	go util.GetByList(urlCnProxy, 1, re.Compile2, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "cn-proxy.com"}
	}
}
