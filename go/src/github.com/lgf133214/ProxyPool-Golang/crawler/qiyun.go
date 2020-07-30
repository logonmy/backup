package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
	"strconv"
)

func qiYunRun() {
	qiYunUrl := "https://www.7yip.cn/free/?page="

	ch := make(chan string, 100)
	list := make([]string, 0, 7)
	for i := range [7][0]int{} {
		list = append(list, qiYunUrl+strconv.Itoa(i+1))
	}

	go util.GetByList(list, 1, re.Compile2, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "www.7yip.cn"}
	}
}
