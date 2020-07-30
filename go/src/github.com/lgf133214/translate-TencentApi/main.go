package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var (
	// 自己去官网开账号，开应用
	appId  = "2135464077"
	appKey = "0SATjdhcApt3EWpt"
)

// 为了排序
type Param struct {
	key, value string
}

func SortParams(p []Param) {
	// 升序
	sort.Slice(p, func(i, j int) bool {
		if p[i].key < p[j].key {
			return true
		}
		return false
	})
}

func ParamsToString(p []Param) string {
	s := ""
	for _, v := range p {
		if v.value == "" {
			continue
		}
		// value 需要进行 url 编码
		s += v.key + "=" + url.QueryEscape(v.value) + "&"
	}
	return s[:len(s)-1]
}

func PostToTansApi(text, target string) (s []byte, ok bool) {
	client := http.Client{}
	params := make([]Param, 0, 10)
	params = append(params, Param{"app_id", appId})
	params = append(params, Param{"source", "auto"})
	params = append(params, Param{"target", target})
	params = append(params, Param{"text", text})

	// 随便啥都行，非空且长度小于32
	params = append(params, Param{"nonce_str", "asbfiuasbhjbcuicg"})
	params = append(params, Param{"time_stamp", fmt.Sprintf("%d", time.Now().Unix())})

	// 取得动态加密标志，获取body
	params = getSign(params)
	form := ParamsToString(params)
	body := bytes.NewBufferString(form)

	request, err := http.NewRequest("POST", "https://api.ai.qq.com/fcgi-bin/nlp/nlp_texttranslate", body)
	if err != nil {
		log.Println(err)
		return
	}

	// 重要的！！！
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}

	defer response.Body.Close()
	s, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	return s, true
}

func getSign(p []Param) []Param {
	s := strings.Builder{}

	SortParams(p)
	s.WriteString(ParamsToString(p))

	s.WriteString("&app_key=" + appKey)

	// MD5
	hash := md5.New()
	hash.Write([]byte(s.String()))
	encodeToString := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))

	p = append(p, Param{"sign", encodeToString})
	return p
}

func TansToZHHandle(c *gin.Context) {
	text, ok := c.GetQuery("text")
	if !ok || text == "" {
		c.JSON(http.StatusOK, map[string]string{"status": "no", "msg": "参数呢"})
		return
	}

	s, ok := PostToTansApi(text, "zh")
	if !ok {
		c.JSON(http.StatusOK, map[string]string{"status": "no", "msg": "不可能出现的错误，多次失败请联系管理员"})
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "ok", "msg": string(s)})
}

func TansToENHandle(c *gin.Context) {
	text, ok := c.GetQuery("text")
	if !ok || text == "" {
		c.JSON(http.StatusOK, map[string]string{"status": "no", "msg": "参数呢"})
		return
	}

	s, ok := PostToTansApi(text, "en")
	if !ok {
		c.JSON(http.StatusOK, map[string]string{"status": "no", "msg": "不可能出现的错误，多次失败请联系管理员"})
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "ok", "msg": string(s)})
}

func AllowControl(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
}

func TlsHandler(c *gin.Context) {
	secureMiddleware := secure.New(secure.Options{
		SSLRedirect: true,
		SSLHost:     ":8203",
	})
	err := secureMiddleware.Process(c.Writer, c.Request)

	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(TlsHandler)
	r.Use(AllowControl)
	{
		r.GET("trans/to/zh", TansToZHHandle)
		r.GET("trans/to/en", TansToENHandle)
	}

	err := r.RunTLS(":8203", "2893797_ligaofeng.top.pem", "2893797_ligaofeng.top.key")
	if err != nil {
		panic(err)
	}
}

/*
本地测试注释掉
r.Use(TlsHandler())
改
err := r.RunTLS(":8203", "*.pem", "*.key")
为
err := r.Run(":8203")
*/
