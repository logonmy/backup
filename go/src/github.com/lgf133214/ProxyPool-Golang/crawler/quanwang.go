package crawler

import (
	"github.com/PuerkitoBio/goquery"

	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/re"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
	"io/ioutil"
	"net/http"
	"strings"
)

func quanWangRun() {
	ch := make(chan string, 100)

	go parseHtml(getHtml(), ch)

	for i := range ch {
		storage.BufferChan <- model.BufferItem{ProxyIp: i, Source: "www.goubanjia.com"}
	}
}

func parseHtml(html string, ch chan<- string) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("www.goubanjia.com " + err.Error())
		close(ch)
		return
	}

	dom.Find("tbody>tr").Each(func(i int, selection *goquery.Selection) {
		selection = selection.Find("td").First()
		s := ""
		val := selection.Find(".port")
		port, ok := val.Attr("class")
		val.Remove()
		if ok {
			port = strings.ReplaceAll(port, "port ", "")
			port = util.GetRealPort(port)
			selection.Children().Each(func(i int, selection *goquery.Selection) {
				attr, ok := selection.Attr("style")
				if ok {
					if !strings.Contains(attr, "display:none") {
						s += selection.Text()
					}
				}
			})
			s += ":" + port
			s = strings.ReplaceAll(s, " ", "")
			if re.Compile1.MatchString(s) {
				ch <- s
			}
		}
	})
	close(ch)
}

func getHtml() string {
	logger.Logger.Debug("visit http://www.goubanjia.com/")
	client := http.Client{}
	request, err := http.NewRequest("GET", "http://www.goubanjia.com/", nil)
	if err != nil {
		log.Println("www.goubanjia.com " + err.Error())
		return ""
	}
	ua := util.GetUA()
	request.Header.Set("User-Agent", ua)

	resp, err := client.Do(request)
	if err != nil {
		log.Println("www.goubanjia.com " + err.Error())
		return ""
	}
	defer resp.Body.Close()

	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("www.goubanjia.com " + err.Error())
		return ""
	}
	return string(all)
}
