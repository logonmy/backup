package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
)

var (
	startUrl = "https://www.xicidaili.com/"
	domain   = []string{"www.xicidaili.com"}
)

// 特别容易503
// 好像后半天就不行
func xiCiRun() {
	ch := make(chan string, 100)

	// depth 2 is enough, change to 3 is too many
	go util.GetByDepth(startUrl, domain, 2, 1, 0, re.Compile2, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "www.xicidaili.com"}
	}
}
