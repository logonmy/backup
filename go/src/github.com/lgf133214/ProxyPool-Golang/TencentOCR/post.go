package TencentOCR

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
可以直接拿去使用
把自己开通好的 id 和 key 换了就可以，需要在应用中添加对应能力
返回的数据是原始字节数组（json格式），根据需要可以直接返回给前端解析
或者在服务端解析
 */

var (
	// 自己去官网开账号，开应用
	appId  = "2135464077"
	appKey = "0SATjdhcApt3EWpt"
	// Tencent's OCR QPS is 2 for me
	OcrLimiter = rate.NewLimiter(2, 2)
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

func GetByteFromUrl(url string) (ret []byte, ok bool) {
	for {
		if ok := OcrLimiter.Allow(); ok {
			break
		}
		time.Sleep(time.Millisecond * 550)
	}
	fileName := "pic" + strconv.FormatInt(time.Now().UnixNano(), 10)
	if ok := downloadPic(url, fileName); !ok {
		return
	}
	defer delFile("./temp/" + fileName)
	return GetByteFromFile("./temp/" + fileName)
}

func GetByteFromFile(filePath string) (ret []byte, ok bool) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	pic, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}

	encodeToString := base64.StdEncoding.EncodeToString(pic)

	return GetByteFromBase64String(encodeToString)
}

func GetByteFromBase64String(encodeToString string) (ret []byte, ok bool) {
	client := http.Client{}
	params := make([]Param, 0, 10)
	params = append(params, Param{"app_id", appId})
	params = append(params, Param{"image", encodeToString})

	// 随便啥都行，非空且长度小于32
	params = append(params, Param{"nonce_str", "asbfiuasbhjbcuicg"})
	params = append(params, Param{"time_stamp", fmt.Sprintf("%d", time.Now().Unix())})

	// 取得动态加密标志，组装body
	params = getSign(params)
	form := ParamsToString(params)
	body := bytes.NewBufferString(form)

	request, err := http.NewRequest("POST", "https://api.ai.qq.com/fcgi-bin/ocr/ocr_generalocr", body)
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
	ret, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	return ret, true
}

// 下载图片
func downloadPic(url, fileName string) (ok bool) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	err = os.Mkdir("./temp", os.ModePerm)
	if !os.IsExist(err) {
		log.Println(err)
		return
	}

	// 后缀不重要，就省去了
	file, err := os.OpenFile("./temp/"+fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	return true
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

func delFile(filePath string) {
	os.Remove(filePath)
}
