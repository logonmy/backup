package crawler

import (
	"fmt"

	"reflect"
	"runtime"
	"time"
)

type f func()

func crawlFunc(f func(), seg time.Duration) func() {
	return func() {
		f()
		funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		logger.Logger.Debug(fmt.Sprintf("%s run over one time", funcName))
		for {
			select {
			case <-time.After(seg):
				f()
				logger.Logger.Debug(fmt.Sprintf("%s run over one time", funcName))
			}
		}
	}
}

func Run() {
	crawlerList := []f{
		//crawlFunc(ipGitHubRun, time.Minute*10),
		//crawlFunc(xiCiRun, time.Minute*10),
		//crawlFunc(quanWangRun, time.Minute*10),
		crawlFunc(miPuRun, time.Minute*10),
		//crawlFunc(data5uRun, time.Minute*10),
		//crawlFunc(coolProxyRun, time.Minute*10),
		//crawlFunc(freeProxyListRun, time.Minute*20),
		crawlFunc(qingTingRun, time.Minute*30),
		//crawlFunc(xiaoHuanRun, time.Minute*30),
		//crawlFunc(proxyListRun, time.Minute*30),
		//crawlFunc(cnProxyRun, time.Minute*30),
		//// site closed?
		////crawlFunc(ip66Run, time.Minute*30),
		//crawlFunc(ip89Run, time.Minute*30),
		//crawlFunc(qiYunRun, time.Hour*1),
		//crawlFunc(ip3366Run, time.Hour*1),
		//crawlFunc(kuaiDaiLiRun, time.Hour*2),
		//crawlFunc(xiaoShuRun, time.Hour*12),
		//crawlFunc(feiZhuRun, time.Hour*12),
		//crawlFunc(atomontersoftRun, time.Hour*12),
		//crawlFunc(proxyIpListRun, time.Hour*12),
		//crawlFunc(proxyListDaily, time.Hour*12),
	}
	for _, i := range crawlerList {
		go i()
	}
}
