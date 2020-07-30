package crawler

import (
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	Run()
	spiders := map[string]int{
		//"www.xicidaili.com":      0,
		//"www.xsdaili.com":        0,
		//"ip.ihuan.me":            0,
		//"www.goubanjia.com":      0,
		//"www.7yip.cn":            0,
		"proxy.horocn.com":       0,
		//"proxy-list.org":         0,
		"proxy.mimvp.com":        0,
		//"www.kuaidaili.com":      0,
		//"ip.jiangxianli.com":     0,
		//"www.ip3366.net":         0,
		//"www.89ip.cn":            0,
		////"www.66ip.cn":            0,
		//"www.feizhuip.com":       0,
		//"www.data5u.com":         0,
		//"www.atomintersoft.com":  0,
		//"www.cool-proxy.net":     0,
		//"free-proxy-list.com":    0,
		//"proxy-ip-list.com":      0,
		//"www.proxylistdaily.net": 0,
		//"cn-proxy.com":           0,
	}

	go func() {
		select {
		case <-time.After(time.Minute * 5):
			close(storage.BufferChan)
		}
	}()
	for i := range storage.BufferChan {
		spiders[i.Source]++
	}
	for i, v := range spiders {
		t.Logf("%s:%d", i, v)
	}
}
