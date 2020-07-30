package network

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/lgf133214/UnblockNeteaseMusic-GO/common"
	"github.com/lgf133214/UnblockNeteaseMusic-GO/utils"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type Netease struct {
	Path     string
	Params   string
	JsonBody map[string]interface{}
}
type ClientRequest struct {
	Method               string
	RemoteUrl            string
	Host                 string
	ForbiddenEncodeQuery bool
	Header               http.Header
	Body                 io.Reader
	Cookies              []*http.Cookie
	Proxy                bool
	ConnectTimeout       time.Duration
}

func Request(clientRequest *ClientRequest) (*http.Response, error) {
	//fmt.Println(clientRequest.RemoteUrl)
	method := clientRequest.Method
	remoteUrl := clientRequest.RemoteUrl
	host := clientRequest.Host
	header := clientRequest.Header
	body := clientRequest.Body
	proxy := clientRequest.Proxy
	cookies := clientRequest.Cookies
	connectTimeout := clientRequest.ConnectTimeout
	if connectTimeout == 0 {
		connectTimeout = 10 * time.Second
	}
	var resp *http.Response
	request, err := http.NewRequest(method, remoteUrl, body)
	if err != nil {
		fmt.Printf("NewRequest fail:%v\n", err)
		return resp, nil
	}
	if !clientRequest.ForbiddenEncodeQuery {
		request.URL.RawQuery = request.URL.Query().Encode()
	}
	c := http.Client{}
	tr := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   connectTimeout,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if len(host) > 0 {
		request.Host = host
		request.Header.Set("host", host)
	}
	if len(request.URL.Scheme) == 0 {
		if request.TLS != nil {
			request.URL.Scheme = "https"
		} else {
			request.URL.Scheme = "http"
		}
	}
	if request.URL.Scheme == "https" || request.TLS != nil {
		if _, ok := common.HostDomain[request.Host]; ok {
			tr.TLSClientConfig = &tls.Config{}
			// verify music.163.com certificate
			tr.TLSClientConfig.ServerName = request.Host //it doesn't contain any IP SANs
			// redirect to music.163.com will need verify self
			tr.TLSClientConfig.InsecureSkipVerify = true
			c.Transport = &tr
		}
	}
	if proxy { //keep headers&cookies for Direct
		if header != nil {
			request.Header = header
		}
		for _, value := range cookies {
			request.AddCookie(value)
		}
	}
	accept := "application/json, text/plain, */*"
	acceptEncoding := "gzip, deflate"
	acceptLanguage := "zh-CN,zh;q=0.9"
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"

	if header != nil {
		accept = header.Get("accept")
		if len(accept) == 0 {
			accept = "application/json, text/plain, */*"
		}
		acceptEncoding = header.Get("accept-encoding")
		if len(acceptEncoding) == 0 {
			acceptEncoding = "gzip, deflate"
		}
		acceptLanguage = header.Get("accept-language")
		if len(acceptLanguage) == 0 {
			acceptLanguage = "zh-CN,zh;q=0.9"
		}
		userAgent = header.Get("user-agent")
		if len(userAgent) == 0 {
			userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"
		}
		Range := header.Get("range")
		if len(Range) > 0 {
			request.Header.Set("range", Range)
		}
	}

	request.Header.Set("accept", accept)
	request.Header.Set("accept-encoding", acceptEncoding)
	request.Header.Set("accept-language", acceptLanguage)
	request.Header.Set("user-agent", userAgent)

	resp, err = c.Do(request)

	if err != nil {
		//fmt.Println(request.Method, request.URL.String(), host)
		fmt.Printf("http.Client.Do fail:%v\n", err)
		return resp, err
	}

	return resp, err

}
func StealResponseBody(resp *http.Response) (io.Reader, error) {
	encode := resp.Header.Get("Content-Encoding")
	enableGzip := false
	if len(encode) > 0 && (strings.Contains(encode, "gzip") || strings.Contains(encode, "deflate")) {
		enableGzip = true
	}
	if enableGzip {
		resp.Header.Del("Content-Encoding")
		body, err := utils.UnGzipV2(resp.Body)
		if err != nil {
			fmt.Println("read  body fail")
			return resp.Body, err
		}
		return body, err
	}
	return resp.Body, nil

}
func GetResponseBody(resp *http.Response, keepBody bool) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body fail")
		return body, err
	}
	resp.Body.Close()
	if keepBody {
		bodyHold := ioutil.NopCloser(bytes.NewBuffer(body))
		resp.Body = bodyHold
	}
	encode := resp.Header.Get("Content-Encoding")
	enableGzip := false
	if len(encode) > 0 && (strings.Contains(encode, "gzip") || strings.Contains(encode, "deflate")) {
		enableGzip = true
	}
	if enableGzip {
		if !keepBody {
			resp.Header.Del("Content-Encoding")
		}
		body, err = utils.UnGzip(body)
		if err != nil {
			fmt.Println("read  body fail")
			return body, err
		}
	}
	return body, err
}
