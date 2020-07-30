package crawler

import (
	"encoding/json"

	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"io/ioutil"
	"net/http"
	"strconv"
)

func coolProxyRun() {
	ch := make(chan string, 100)
	go parseJson(ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "www.cool-proxy.net"}
	}
}

func parseJson(ch chan string) {
	url := "http://www.cool-proxy.net/proxies.json"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	ret := make([]interface{}, 0, 100)
	err = json.Unmarshal(all, &ret)
	if err != nil {
		log.Println(err)
		return
	}
	for _, i := range ret {
		proxy := ""
		if data, ok := i.(map[string]interface{}); ok {
			if ip, ok := data["ip"]; ok {
				if _, ok := ip.(string); ok {
					proxy += ip.(string)
				}
				if port, ok := data["port"]; ok {
					if _, ok := port.(float64); ok {
						proxy += ":" + strconv.Itoa(int(port.(float64)))
					}
				}
			}
		}
		ch <- proxy
	}
	close(ch)
}
