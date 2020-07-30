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
	urlsMiPu = []string{
		"https://proxy.mimvp.com/freesecret?proxy=in_hp",
		"https://proxy.mimvp.com/freesecret?proxy=in_socks",
		"https://proxy.mimvp.com/freesole?proxy=in_hp",
		"https://proxy.mimvp.com/freesole?proxy=in_socks",
		"https://proxy.mimvp.com/freeopen?proxy=in_hp",
		"https://proxy.mimvp.com/freeopen?proxy=in_tp",
		"https://proxy.mimvp.com/freeopen?proxy=in_socks",
		"https://proxy.mimvp.com/freeopen?proxy=out_hp",
		"https://proxy.mimvp.com/freeopen?proxy=out_tp",
		"https://proxy.mimvp.com/freeopen?proxy=out_socks",
	}
	compileCustom = regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)(?s:.*?)<img src=(.*?) />`)
)

// 迭代很快了
func miPuRun() {
	ch := make(chan string, 100)

	go util.GetByList(urlsMiPu, 1, compileCustom, ch)

	for i := range ch {
		data := strings.Split(i, ":")
		imgUrl := "https://proxy.mimvp.com" + data[1]
		port := TencentOCR.GetByteFromUrl(imgUrl)
		if port != "" {
			if util.FilterPort(port) {
				storage.BufferChan <- model.BufferItem{ProxyIp: data[0] + ":" + port, Source: "proxy.mimvp.com"}
			}
		}
	}
}
