package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
)

var (
	urlXiaoShu = "http://www.xsdaili.com/"
)

func xiaoShuRun() {
	ch := make(chan string, 100)

	go util.GetByDepth(urlXiaoShu, []string{"www.xsdaili.com"}, 2, 0, 5, re.Compile1, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "www.xsdaili.com"}
	}
}
