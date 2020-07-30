package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lgf133214/wendaku/log_"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	regStr = "\\s{2,}"              //两个及两个以上空格的正则表达式
	reg, _ = regexp.Compile(regStr) //编译正则表达式

	Urls    []*url.URL
	proxies []string
	index   uint32

	mu sync.RWMutex
)

func init() {
	GetProxies()
	go func() {
		for {
			select {
			case <-time.After(time.Minute * 3):
				GetProxies()
			}
		}
	}()
}

func DeleteExtraSpace(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "\n", " ", -1)     //替换\n为空格
	s1 = strings.Replace(s1, " ", " ", -1)      //替换tab为空格
	s2 := make([]byte, len(s1))                 //定义字符数组切片
	copy(s2, s1)                                //将字符串复制到切片
	spcIndex := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spcIndex) > 0 {                     //找到适配项
		s2 = append(s2[:spcIndex[0]+1], s2[spcIndex[1]:]...) //删除多余空格
		spcIndex = reg.FindStringIndex(string(s2))           //继续在字符串中搜索
	}
	return strings.TrimSpace(string(s2))
}

func ProxyFunc(pr *http.Request) (*url.URL, error) {
	mu.RLock()
	u := Urls[index%uint32(len(Urls))]
	mu.RUnlock()

	atomic.AddUint32(&index, 1)
	ctx := context.WithValue(pr.Context(), colly.ProxyURLKey, u.String())
	*pr = *pr.WithContext(ctx)
	return u, nil
}

func GetProxies() {

	items, ok := GetSuccessItems()
	if !ok {
		log_.ErrWrite(errors.New("代理IP获取失败"))
		return
	}

	mu.Lock()
	proxies = make([]string, 0, len(items))
	for _, i := range items {
		proxies = append(proxies, "https://"+i.ProxyIp)
	}

	Urls = make([]*url.URL, len(proxies))
	for i, u := range proxies {
		parsedU, err := url.Parse(u)
		if err != nil {
			log_.ErrWrite(err)
			continue
		}
		Urls[i] = parsedU
	}

	mu.Unlock()
}

func GenIpAddr() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}
