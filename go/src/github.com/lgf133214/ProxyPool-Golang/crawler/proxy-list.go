package crawler

import (
	"encoding/base64"

	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
	"regexp"
	"strconv"
)

var (
	proxyListUrl = "https://proxy-list.org/chinese/index.php?p="
	proxyListRe  = regexp.MustCompile(`>Proxy\('(.*?)'\)<`)
)

func proxyListRun() {
	ch := make(chan string, 100)
	list := make([]string, 0, 10)
	for i := range [10][0]int{} {
		list = append(list, proxyListUrl+strconv.Itoa(i+1))
	}

	go util.GetByList(list, 1, proxyListRe, ch)

	for i := range ch {
		s, err := base64.StdEncoding.DecodeString(i)
		if err != nil {
			log.Println("proxy-list.org " + err.Error())
			continue
		}
		storage.BufferChan <- model.BufferItem{ProxyIp: string(s), Source: "proxy-list.org"}
	}
}
