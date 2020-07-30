package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/TencentOCR"
	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
	"regexp"
	"strings"
)

var (
	urlQingTing = "https://proxy.horocn.com/free-proxy.html"
	qingTingRe  = regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)(?s:.*?)base64,(.*?)">`)
)

func qingTingRun() {
	ch := make(chan string, 100)

	go util.GetByDepth(urlQingTing, []string{"proxy.horocn.com"}, 4, 1, 0, qingTingRe, ch)

	for i := range ch {
		data := strings.Split(i, ":")
		port := TencentOCR.GetByteFromBase64String(data[1])
		if port != "" {
			if util.FilterPort(port) {
				storage.BufferChan <- model.BufferItem{ProxyIp: data[0] + ":" + port, Source: "proxy.horocn.com"}
			}
		}
	}

}
