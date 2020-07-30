package server

import (
	"encoding/json"

	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"strings"
)

var (
	httpLimiter  = rate.NewLimiter(800, 1500)
	ipLimiter    = NewIPRateLimiter(1, 2)
)

func init() {
	http.HandleFunc("/", wrapHandler)
}

func Run() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("http server run error")
		return
	}
}

func wrapHandler(w http.ResponseWriter, r *http.Request) {

	ip := strings.Split(r.RemoteAddr, ":")[0]
	limiter := ipLimiter.GetLimiter(ip)
	ok := limiter.Allow()
	if !ok {
		w.Write([]byte(`{"data":"no results"}`))
		return
	}
	ok = httpLimiter.Allow()
	if !ok {
		w.Write([]byte(`{"data":"no results"}`))
		return
	}

	paras := r.URL.Query()

	logger.Logger.Debug("request from " + r.RemoteAddr)
	var skip int64
	if util.GetPara(paras, "skip") != "" {
		i, err := strconv.ParseInt(util.GetPara(paras, "skip"), 10, 64)
		if err == nil {
			skip = i
		}
	}

	count, items, ok := storage.GetSuccessItems(skip, paras)
	if !ok {
		w.Write([]byte(`{"data":"no results"}`))
		return
	}
	var ret = make(map[string]interface{})
	ret["status"] = "ok"
	ret["count"] = count
	var data []interface{}
	for _, i := range items {
		tmp := make(map[string]interface{})
		tmp["proxy_ip"] = i.ProxyIp
		tmp["http"] = i.Http
		tmp["https"] = i.Https
		tmp["google"] = i.Google
		tmp["socks5"] = i.Socks5
		tmp["verify_time"] = i.VerifyTime
		tmp["anonymous"] = i.Anonymous
		data = append(data, tmp)
	}
	ret["data"] = data
	bytes, err := json.Marshal(ret)
	if err != nil {
		log.Println(err)
		w.Write([]byte(`{"data":"no results"}`))
		return
	}
	w.Write(bytes)
}
