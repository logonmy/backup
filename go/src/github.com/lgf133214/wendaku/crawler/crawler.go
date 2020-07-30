package crawler

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/jinzhu/gorm"
	"github.com/lgf133214/wendaku/log_"
	"github.com/lgf133214/wendaku/util"
	"github.com/zolamk/colly-mongo-storage/colly/mongo"
	"log"
	"strings"
)

type ItemHasAnswer struct {
	Id       int    `gorm:"primary_key"`
	Question string `gorm:"unique_index;type:varchar(1023);not null"`
	Answer   string `gorm:"index;type:varchar(1023);not null"`
}

type ItemHasNoAnswer struct {
	Id       int    `gorm:"primary_key"`
	Question string `gorm:"unique_index;type:varchar(1023);not null"`
	Options  string `gorm:"index;type:varchar(1023);not null"`
}

type ItemOnlyHasQuestion struct {
	Id       int    `gorm:"primary_key"`
	Question string `gorm:"unique_index;type:varchar(1023);not null"`
}

var (
	connInfo = "root:58601a7f241645fe@/wen_da_ku?charset=utf8&parseTime=True&loc=Local"
)

func init() {
	var db, err = gorm.Open("mysql", connInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Set("gorm:wen_da_ku", "ENGINE=InnoDB").AutoMigrate(&ItemHasAnswer{}, &ItemHasNoAnswer{}, &ItemOnlyHasQuestion{})
}

func Run() {
	c := colly.NewCollector(func(c *colly.Collector) {
		c.AllowURLRevisit = false
		c.AllowedDomains = []string{"www.asklib.com"}
		extensions.RandomUserAgent(c)
		extensions.Referer(c)
		c.SetProxyFunc(util.ProxyFunc)
		c.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: len(util.Urls),
		})
	})
	c2 := colly.NewCollector(func(c *colly.Collector) {
		c.AllowURLRevisit = false
		c.AllowedDomains = []string{"www.asklib.com"}
		extensions.RandomUserAgent(c)
		extensions.Referer(c)
		c.SetProxyFunc(util.ProxyFunc)
		c.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: len(util.Urls),
		})
	})

	storage1 := &mongo.Storage{
		Database: "others",
		URI:      "mongodb://127.0.0.1:27017",
	}
	storage2 := &mongo.Storage{
		Database: "views",
		URI:      "mongodb://127.0.0.1:27017",
	}

	c.SetStorage(storage1)
	c2.SetStorage(storage2)

	c2.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Forwarded-For", util.GenIpAddr())
	})
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Forwarded-For", util.GenIpAddr())
	})

	// 提取标签
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.Contains(href, ".html") {
			if strings.Contains(href, "/view/") {
				c2.Visit(e.Request.AbsoluteURL(href))
			} else {
				e.Request.Visit(e.Request.AbsoluteURL(href))
			}
		}
	})

	c2.OnHTML(".content.clear", func(e *colly.HTMLElement) {
		// 获取题目
		tm := e.ChildText("h1")
		if tm == "" {
			log_.ErrWrite(errors.New(fmt.Sprintf("题目为空，%s", e.Request.AbsoluteURL(""))))
			return
		}
		tm = util.DeleteExtraSpace(tm)

		// 获取选项
		xx, err := e.DOM.Find("h2+p").Html()
		if err != nil {
			log_.ErrWrite(errors.New(fmt.Sprintf("选项解析失败，%s", e.Request.AbsoluteURL(""))))
			return
		}
		xx = util.DeleteExtraSpace(xx)

		// 尝试获取答案
		da := e.ChildText(".listbg")
		da = util.DeleteExtraSpace(da)

		if da == "参考答案： 查看" || da == "" {
			da = ""
		} else {
			if strings.Contains(da, "答案：") {
				da = strings.TrimSpace(strings.Split(da, "答案：")[1])
			} else if !strings.Contains(da, "关键词：") {
				log_.ErrWrite(errors.New(fmt.Sprintf("答案捕获错误，%s 题目：%s 答案：%s 选项：%s",
					e.Request.AbsoluteURL(""), tm, da, xx)))
			}
		}

		log.Printf("url：%s 题目：%s 答案：%s 选项：%s", e.Request.AbsoluteURL(""), tm, da, xx)

		// 选项为空
		if xx == "" {
			if da == "" {
				// 入题目库
				i := ItemOnlyHasQuestion{Question: tm}
				store(i)
			} else {
				// 答案不为空
				i := ItemHasAnswer{Question: tm, Answer: da}
				store(i)
			}
		} else {
			if da == "" {
				i := ItemHasNoAnswer{Question: tm, Options: xx}
				store(i)
			} else {
				// 根据单选、多选信息分辨
				xxs := strings.Split(xx, "<br/>")
				i := ItemHasAnswer{Question: tm}
				if strings.Contains(tm, "单选") {
					for _, j := range xxs {
						j = strings.TrimSpace(j)
						if strings.Contains(j, da) {
							i.Answer = j
						}
					}
					if i.Answer == "" {
						log_.ErrWrite(errors.New("匹配不到答案，" + e.Request.AbsoluteURL("") + tm + "--" + da))
						return
					}
				} else if strings.Contains(tm, "多选") || strings.Contains(tm, "不定项选择") {
					das := strings.Split(da, ",")
					for _, j := range xxs {
						j = strings.TrimSpace(j)
						for _, x := range das {
							x = strings.TrimSpace(x)
							if strings.Contains(j, x) {
								i.Answer += " " + j
							}
						}
					}
					if i.Answer == "" {
						log_.ErrWrite(errors.New("匹配不到答案，" + e.Request.AbsoluteURL("") + tm + "--" + da))
						return
					}
				} else if strings.Contains(tm, "判断") {
					i.Answer = da
				} else {
					log_.ErrWrite(errors.New("未知类型，" + e.Request.AbsoluteURL("") + tm + "--" + da))
					return
				}
				store(i)
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log_.ErrWrite(err)
		r.Request.Visit(r.Request.URL.String())
	})

	c2.OnError(func(r *colly.Response, err error) {
		log_.ErrWrite(err)
		r.Request.Visit(r.Request.URL.String())
	})

	c.Visit("https://www.asklib.com")

	log.Println("crawler end")
}

func store(i interface{}) {
	switch i.(type) {
	case ItemHasAnswer:
		var db, err = gorm.Open("mysql", connInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		x := ItemHasAnswer{}
		tmp := i.(ItemHasAnswer)
		db.Find(&x, "question=?", tmp.Question)

		if x.Id == 0 {
			db.Create(&tmp)
		}
	case ItemHasNoAnswer:
		var db, err = gorm.Open("mysql", connInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		x := ItemHasNoAnswer{}
		tmp := i.(ItemHasNoAnswer)
		db.Find(&x, "question=?", tmp.Question)

		if x.Id == 0 {
			db.Create(&tmp)
		}
	case ItemOnlyHasQuestion:
		var db, err = gorm.Open("mysql", connInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		x := ItemOnlyHasQuestion{}
		tmp := i.(ItemOnlyHasQuestion)
		db.Find(&x, "question=?", tmp.Question)

		if x.Id == 0 {
			db.Create(&tmp)
		}
	default:
		log.Printf("%T，类型错误，不包含该类型！", i)
	}
}
