package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
)


func feiZhuRun() {
	urlFeiZhu := "https://www.feizhuip.com/News-newsList-catid-8.html"

	ch := make(chan string, 100)
	go util.GetByDepth(urlFeiZhu, []string{"www.feizhuip.com"}, 2, 0, 5, re.Compile2, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "www.feizhuip.com"}
	}
}
