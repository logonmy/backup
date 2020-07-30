package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
	"strconv"
)

func ip3366Run() {
	url3366T1 := "http://www.ip3366.net/free/?stype=1&page="
	url3366T2 := "http://www.ip3366.net/free/?stype=2&page="

	ch := make(chan string, 100)
	list := make([]string, 0, 10)
	for i := range [6][0]int{} {
		list = append(list, url3366T1+strconv.Itoa(i+1))
		list = append(list, url3366T2+strconv.Itoa(i+1))
	}

	go util.GetByList(list, 1, re.Compile2, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "www.ip3366.net"}
	}
}
