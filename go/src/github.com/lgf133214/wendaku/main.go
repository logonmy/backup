package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	//crawler.Run()

	proxy_ := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + "46.151.145.4:53281")
	}
	client := http.Client{Transport: &http.Transport{Proxy: proxy_}}
	resp, err := client.Get("http://httpbin.org/get")
	if err != nil {
		panic(err)
	}
	all, _ := ioutil.ReadAll(resp.Body)
	log.Printf("%s", all)
}


