package util

import (
	"encoding/json"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var LocalIp = GetLocalIp()

// every 12 hour check the local ip address by default
// of course, you can cancel it if your machine is fixed ip
func init() {
	go func() {
		for {
			select {
			case <-time.After(time.Hour * 12):
				LocalIp = GetLocalIp()
			}
		}
	}()
}

// 0 Transparent 1 Distorting 2 Anonymous
// return true means that the proxy protocol is http
func HttpAnonymous(proxyIp string) (int, bool) {
	proxy_ := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + proxyIp)
	}

	data, ok := RequestHttpBinByHttp(proxy_)
	if !ok {
		return 0, false
	}

	// get origin
	origin, ok := ParseOrigin(data)
	if !ok {
		return 0, false
	}

	// the anonymous proxy will not show your local ip
	if strings.Contains(origin, ",") {
		if strings.Contains(origin, LocalIp) {
			return 1, true
		} else {
			return 2, true
		}
	} else {
		if strings.Contains(origin, LocalIp) {
			return 0, true
		} else {
			return 2, true
		}
	}
}

// get a string like "ip1, ip2...", you can view detail pattern on http://httpbin.org/get
func ParseOrigin(data []byte) (string, bool) {
	origin := make(map[string]interface{})
	err := json.Unmarshal(data, &origin)
	if err != nil {
		return "", false
	}

	val, ok := origin["origin"]
	if !ok {
		return "", false
	}

	s, ok := val.(string)
	if !ok {
		return "", false
	}
	return s, true
}

// send request to http://httpbing.org/get
func RequestHttpBinByHttp(proxy func(_ *http.Request) (*url.URL, error)) ([]byte, bool) {
	client := http.Client{
		Timeout: time.Second * 20,
	}
	if proxy != nil {
		client.Transport = &http.Transport{Proxy: proxy}
	}
	request, err := http.NewRequest("GET", "http://httpbin.org/get", nil)
	if err != nil {
		return nil, false
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, false
	}
	defer response.Body.Close()
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, false
	}
	return all, true
}

// send request to https://httpbin.org/get
func RequestHttpBinByHttps(proxyIp string) bool {
	proxy_ := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + proxyIp)
	}
	client := http.Client{
		Timeout: time.Second * 20,
	}
	if proxy_ != nil {
		client.Transport = &http.Transport{Proxy: proxy_}
	}
	request, err := http.NewRequest("GET", "https://httpbin.org/get", nil)
	if err != nil {
		return false
	}
	_, err = client.Do(request)
	if err != nil {
		return false
	}
	return true
}

func RequestGoogleByHttps(proxyIp string) bool {
	proxy_ := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + proxyIp)
	}
	client := http.Client{
		Timeout: time.Second * 20,
	}
	if proxy_ != nil {
		client.Transport = &http.Transport{Proxy: proxy_}
	}
	request, err := http.NewRequest("GET", "https://www.google.com/robots.txt", nil)
	if err != nil {
		return false
	}

	// if err is not nil, the request must be bad
	_, err = client.Do(request)
	if err != nil {
		return false
	}
	return true
}

func ValSocks5(proxyIp string) bool {
	dialer, err := proxy.SOCKS5("tcp", proxyIp, nil, proxy.Direct)
	if err != nil {
		return false
	}

	httpClient := &http.Client{Transport: &http.Transport{Dial: dialer.Dial}, Timeout: time.Second * 20}
	if _, err := httpClient.Get("https://httpbin.org/get"); err != nil {
		return false
	} else {
		return true
	}
}

func RequestGoogleBySocks5(proxyIp string) bool {
	dialer, err := proxy.SOCKS5("tcp", proxyIp, nil, proxy.Direct)
	if err != nil {
		return false
	}

	httpClient := &http.Client{Transport: &http.Transport{Dial: dialer.Dial}, Timeout: time.Second * 20}
	if _, err := httpClient.Get("https://www.google.com/robots.txt"); err != nil {
		return false
	} else {
		return true
	}
}
