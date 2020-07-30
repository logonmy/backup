package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
)

// 最多2800，足够
func ip89Run() {
	ch := make(chan string, 100)
	go util.GetByList([]string{"http://www.89ip.cn/tqdl.html?num=2800"}, 1, re.Compile1, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "www.89ip.cn"}
	}
}
