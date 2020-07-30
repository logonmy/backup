package util

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"

	"github.com/lgf133214/ProxyPool-Golang/re"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// get User-Agent function come from colly extensions
var uaGens = []func() string{
	genFirefoxUA,
	genChromeUA,
}

var ffVersions = []float32{
	58.0,
	57.0,
	56.0,
	52.0,
	48.0,
	40.0,
	35.0,
}

var chromeVersions = []string{
	"65.0.3325.146",
	"64.0.3282.0",
	"41.0.2228.0",
	"40.0.2214.93",
	"37.0.2062.124",
}

var osStrings = []string{
	"Macintosh; Intel Mac OS X 10_10",
	"Windows NT 10.0",
	"Windows NT 5.1",
	"Windows NT 6.1; WOW64",
	"Windows NT 6.1; Win64; x64",
	"X11; Linux x86_64",
}

func genFirefoxUA() string {
	version := ffVersions[rand.Intn(len(ffVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%.1f) Gecko/20100101 Firefox/%.1f", os, version, version)
}

func genChromeUA() string {
	version := chromeVersions[rand.Intn(len(chromeVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", os, version)
}

func GetUA() string {
	return uaGens[rand.Intn(len(uaGens))]()
}

// recover unexpected err, instead of exit
func RecoverFunc() {
	if err := recover(); err != nil {
		log.Println(fmt.Sprintf("%s", err))
		return
	}
}

// When depth is 0, request will forever do until the end.
// When depth is 1, it will only request the start page.
// compile is the specific expression you need.
// baseDelay is the time will wait before next request, and
// the final time is baseDelay add a random time between 0-2
// limitPages is num of the response will get
func GetByDepth(startUrl string, domain []string, depth, baseDelay, limitPages int, compile *regexp.Regexp, ch chan<- string) {
	defer close(ch)
	defer RecoverFunc()
	visitNum := 0
	c1 := colly.NewCollector(func(c *colly.Collector) {
		c.AllowURLRevisit = false
		c.AllowedDomains = domain
		extensions.RandomUserAgent(c)
		extensions.Referer(c)
		c.IgnoreRobotsTxt = true
		c.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			RandomDelay: time.Second * 2,
			Delay:       time.Second * time.Duration(baseDelay),
		})
		c.MaxDepth = depth
	})
	c1.OnResponse(func(r *colly.Response) {
		logger.Logger.Debug(r.Request.AbsoluteURL(""))
		data := string(r.Body)
		stringSubmatch := compile.FindAllStringSubmatch(data, -1)

		if stringSubmatch != nil {
			for _, i := range stringSubmatch {
				ch <- strings.Join(i[1:], ":")
			}
			visitNum++
		}
	})
	c1.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if limitPages > 0 && visitNum >= limitPages {
			return
		}
		// 这里的错误不能panic，会返回 Max depth limit reached，麻烦，别管了
		e.Request.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})
	c1.OnError(func(r *colly.Response, err error) {
		log.Println(r.Request.AbsoluteURL("") + " " + err.Error())
	})

	c1.Visit(startUrl)
}

func GetByList(urls []string, baseDelay int, compile *regexp.Regexp, ch chan<- string) {
	defer RecoverFunc()
	defer close(ch)
	c1 := colly.NewCollector(func(c *colly.Collector) {
		c.AllowURLRevisit = false
		extensions.RandomUserAgent(c)
		extensions.Referer(c)
		c.IgnoreRobotsTxt = true
		c.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			RandomDelay: time.Second * 2,
			Delay:       time.Second * time.Duration(baseDelay),
		})
	})
	c1.OnResponse(func(r *colly.Response) {
		logger.Logger.Debug(r.Request.AbsoluteURL(""))
		data := string(r.Body)
		stringSubmatch := compile.FindAllStringSubmatch(data, -1)

		if stringSubmatch != nil {
			for _, i := range stringSubmatch {
				ch <- strings.Join(i[1:], ":")
			}
		}
	})
	c1.OnError(func(r *colly.Response, err error) {
		log.Println(r.Request.AbsoluteURL("") + " " + err.Error())
	})

	for _, i := range urls {
		c1.Visit(i)
	}
}

func GetRealPort(s string) string {
	key := "ABCDEFGHIZ"
	port := 0
	for _, i := range []rune(s) {
		port *= 10
		port += strings.Index(key, string(i))
	}
	port = port >> 3
	return strconv.Itoa(port)
}

func GetLocalIp() string {
	data, ok := RequestHttpBinByHttp(nil)
	if !ok {
		log.Println("你真惨呐，这事都给你碰上了，肯定是你的问题")
		return ""
	}
	origin, ok := ParseOrigin(data)
	return origin
}

func FilterPort(port string) bool {
	_, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		return false
	}
	return true
}

func FilterProxyIp(proxyIp string) bool {
	return re.Compile1.MatchString(proxyIp)
}

func GetPara(data map[string][]string, para string) string {
	p, ok := data[para]
	if ok {
		if len(p) >= 1 {
			return p[0]
		}
	}
	return ""
}
