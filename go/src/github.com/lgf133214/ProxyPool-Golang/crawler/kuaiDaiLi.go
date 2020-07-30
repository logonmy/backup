package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
	"strconv"
)

func kuaiDaiLiRun() {
	urlKuaiDaiLi1 := "https://www.kuaidaili.com/free/inha/"
	urlKuaiDaiLi2 := "https://www.kuaidaili.com/free/intr/"

	ch := make(chan string, 100)
	list := make([]string, 0, 14)
	for i := range [7][0]int{} {
		list = append(list, urlKuaiDaiLi1+strconv.Itoa(i+1))
		list = append(list, urlKuaiDaiLi2+strconv.Itoa(i+1))
	}

	go util.GetByList(list, 1, re.Compile2, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "www.kuaidaili.com"}
	}
}
