package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
	"regexp"
	"strings"
)

var data5uRe = regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)(?s:.*?)port (\w+)`)

func data5uRun() {
	urlData5u := []string{"http://www.data5u.com/"}

	ch := make(chan string, 100)
	go util.GetByList(urlData5u, 1, data5uRe, ch)
	for i := range ch {
		data := strings.Split(i, ":")
		port := util.GetRealPort(data[1])
		storage.BufferChan <- model.BufferItem{ProxyIp: data[0] + ":" + port, Source: "www.data5u.com"}
	}
}
