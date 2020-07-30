package crawler

import (
	"encoding/json"
	"fmt"

	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"io/ioutil"
	"net/http"
	"time"
)


func ipGitHubRun() {
	urlGitHub := "https://ip.jiangxianli.com/api/proxy_ips?page=1"

	ch := make(chan string, 100)
	go get(urlGitHub, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "ip.jiangxianli.com"}
	}
}

func get(url string, ch chan<- string) {
	logger.Logger.Debug("visiting ip.jiangxianli.com, project from github")
	resp, err := http.Get(url)
	if err != nil {
		log.Println("ip.jiangxianli.com  " + err.Error())
		close(ch)
		return
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ip.jiangxianli.com " + err.Error())
		close(ch)
		return
	}

	stringSubmatch := re.Compile2.FindAllStringSubmatch(string(all), -1)
	if stringSubmatch == nil {
		close(ch)
		return
	}
	for _, i := range stringSubmatch {
		ch <- fmt.Sprintf("%s:%s", i[1], i[2])
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(all, &data)
	if err != nil {
		log.Println("ip.jiangxianli.com " + err.Error())
		close(ch)
		return
	}
	tmp, ok := data["data"].(map[string]interface{})
	if !ok {
		close(ch)
		return
	}
	next, ok := tmp["next_page_url"].(string)
	if !ok {
		close(ch)
		return
	}
	time.Sleep(time.Second)
	get(next, ch)
}
