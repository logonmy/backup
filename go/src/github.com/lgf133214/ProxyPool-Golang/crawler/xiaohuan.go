package crawler

import (
	"bytes"
	"fmt"

	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	keyUrl     = "https://ip.ihuan.me/mouse.do"
	refererUrl = "https://ip.ihuan.me/ti.html"
	regCustom  = regexp.MustCompile(`val\(\"(\w+)\"\);`)
	ua         = ""
)

func xiaoHuanRun() {
	ch := make(chan string, 100)
	ua = util.GetUA()
	key := getKey()

	go post(key, 1500, ch)
	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "ip.ihuan.me"}
	}
}

// 3000一次最大，足够了
func post(key string, num int, ch chan<- string) {
	defer close(ch)

	logger.Logger.Debug("request https://ip.ihuan.me/tqdl.html")
	url := "https://ip.ihuan.me/tqdl.html"

	client := http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewBufferString(fmt.Sprintf("key=%s&num=%d&sort=1",
		key, num)))
	if err != nil {
		log.Println("ip.ihuan.me " + err.Error())
		return
	}

	request.Header.Set("Referer", refererUrl)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", ua)

	response, err := client.Do(request)
	if err != nil {
		log.Println("ip.ihuan.me " + err.Error())
		return
	}
	defer response.Body.Close()

	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("ip.ihuan.me " + err.Error())
		return
	}

	stringSubmatch := re.Compile1.FindAllStringSubmatch(string(all), -1)
	if stringSubmatch == nil {
		log.Println("ip.ihuan.me can not match the ips")
		return
	}

	for _, i := range stringSubmatch {
		ch <- i[0]
	}
}

func getKey() string {
	logger.Logger.Debug("ihuan.me getKey...")
	client := http.Client{}
	// 先要获取cookie
	r1, err := http.NewRequest("GET", refererUrl, nil)
	if err != nil {
		log.Println("ip.ihuan.me key " + err.Error())
		return ""
	}
	// User-Agent 前后必须一致，cookie根据这个给，waste about 1 hours，淦
	r1.Header.Set("User-Agent", ua)

	resp, err := client.Do(r1)
	if err != nil {
		log.Println("ip.ihuan.me key " + err.Error())
		return ""
	}

	request, err := http.NewRequest("GET", keyUrl, nil)
	if err != nil {
		log.Println("ip.ihuan.me key " + err.Error())
		return ""
	}
	for _, i := range resp.Cookies() {
		request.AddCookie(i)
	}

	request.Header.Set("User-Agent", ua)
	request.Header.Set("Referer", refererUrl)

	response, err := client.Do(request)
	if err != nil {
		log.Println("ip.ihuan.me key " + err.Error())
		return ""
	}
	defer response.Body.Close()

	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("ip.ihuan.me key " + err.Error())
		return ""
	}
	stringSubmatch := regCustom.FindStringSubmatch(string(all))
	if stringSubmatch == nil {
		log.Println("ip.ihuan.me can not match the key")
		return ""
	}

	return stringSubmatch[1]
}
