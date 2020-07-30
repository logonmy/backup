package scheduler

import (
	"fmt"
	"github.com/lgf133214/ProxyPool-Golang/crawler"

	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"github.com/lgf133214/ProxyPool-Golang/util"
	"reflect"
	"runtime"
	"time"
)

var (
	worker = make(chan struct{}, 3000)
)

func Run() {
	go crawler.Run()
	go storeBuffer()
	go validator()
}

func storeBuffer() {
	go func() {
		for item := range storage.BufferChan {
			storage.StoreToBufferDB(item)
		}
	}()
}

func validator() {
	go validateBufferItem()
	go validateSuccessItem()
	go validateUnableItem()
}

func validateSuccessItem() {
	go func() {
		funcName := runtime.FuncForPC(reflect.ValueOf(validateSuccessItem).Pointer()).Name()
		for {
			select {
			case <-time.After(time.Minute * 5):
				logger.Logger.Debug(fmt.Sprintf("%s has start at %s, please check if there is an end occur", funcName, time.Now().Format("2006-01-02 15:04:05")))
				// pull items into chan
				for {
					if !storage.GetSuccessItem() {
						time.Sleep(time.Second * 30)
					} else {
						break
					}
				}
				logger.Logger.Debug(funcName + " end")
			}
		}
	}()
	go func() {
		for i := range storage.SuccessValChan {
			SuccessWorkerDo(i)
		}
	}()
}

func validateBufferItem() {
	go func() {
		funcName := runtime.FuncForPC(reflect.ValueOf(validateBufferItem).Pointer()).Name()
		for {
			select {
			case <-time.After(time.Minute * 5):
				logger.Logger.Debug(fmt.Sprintf("%s has start at %s, please check if there is an end occur", funcName, time.Now().Format("2006-01-02 15:04:05")))
				// pull to buffer chan
				for {
					if !storage.GetBufferItem() {
						time.Sleep(time.Second * 30)
					} else {
						break
					}
				}
				logger.Logger.Debug(funcName + " end")
			}
		}
	}()
	go func() {
		for i := range storage.BufferValChan {
			bufferWorkerDo(i)
		}
	}()
}

func validateUnableItem() {
	go func() {
		funcName := runtime.FuncForPC(reflect.ValueOf(validateUnableItem).Pointer()).Name()
		for {
			select {
			case <-time.After(time.Minute * 5):
				logger.Logger.Debug(fmt.Sprintf("%s has start at %s, please check if there is an end occur", funcName, time.Now().Format("2006-01-02 15:04:05")))
				// pull to unable chan
				for {
					if !storage.GetUnableItem() {
						time.Sleep(time.Second * 30)
					} else {
						break
					}
				}
				logger.Logger.Debug(funcName + " end")
			}
		}
	}()
	go func() {
		for i := range storage.UnableValChan {
			unableWorkerDo(i)
		}
	}()
}

func SuccessWorkerDo(i model.SuccessItem) {
	worker <- struct{}{}
	go func() {
		defer func() { <-worker }()
		anonymous, http := util.HttpAnonymous(i.ProxyIp)
		https := util.RequestHttpBinByHttps(i.ProxyIp)
		if !http && !https {
			if util.ValSocks5(i.ProxyIp) {
				i.Socks5 = true
				i.VerifyTime = time.Now().Unix()
				i.Google = util.RequestGoogleBySocks5(i.ProxyIp)
				// update
				storage.StoreToSuccessDB(i)
				return
			}
			// remove
			storage.RemoveSuccessItem(i.ProxyIp)
			storage.StoreToUnableDB(model.UnableItem{ProxyIp: i.ProxyIp, Source: i.Source})
			return
		}
		i.Anonymous = anonymous
		i.VerifyTime = time.Now().Unix()
		i.Http = true
		i.Https = true

		if https {
			i.Google = util.RequestGoogleByHttps(i.ProxyIp)
		}

		// update
		storage.StoreToSuccessDB(i)
		// remove
		storage.RemoveSuccessItem(i.ProxyIp)
	}()
}

func bufferWorkerDo(i model.BufferItem) {
	worker <- struct{}{}
	go func() {
		defer func() { <-worker }()
		anonymous, http := util.HttpAnonymous(i.ProxyIp)
		https := util.RequestHttpBinByHttps(i.ProxyIp)
		if !http && !https {
			if util.ValSocks5(i.ProxyIp) {
				n := model.SuccessItem{ProxyIp: i.ProxyIp, Source: i.Source}
				n.Socks5 = true
				n.VerifyTime = time.Now().Unix()
				n.Google = util.RequestGoogleBySocks5(i.ProxyIp)
				// insert
				storage.StoreToSuccessDB(n)
			}
			// remove
			storage.RemoveBufferItem(i.ProxyIp)
			return
		}

		n := model.SuccessItem{ProxyIp: i.ProxyIp, Source: i.Source}
		n.VerifyTime = time.Now().Unix()
		n.Http = http
		n.Https = https
		n.Anonymous = anonymous
		if https {
			n.Google = util.RequestGoogleByHttps(i.ProxyIp)
		}
		// insert
		storage.StoreToSuccessDB(n)
		// remove
		storage.RemoveBufferItem(i.ProxyIp)
	}()
}

func unableWorkerDo(i model.UnableItem) {
	worker <- struct{}{}
	go func() {
		defer func() { <-worker }()
		anonymous, http := util.HttpAnonymous(i.ProxyIp)
		https := util.RequestHttpBinByHttps(i.ProxyIp)
		if !http && !https {
			if util.ValSocks5(i.ProxyIp) {
				n := model.SuccessItem{ProxyIp: i.ProxyIp, Source: i.Source}
				n.Socks5 = true
				n.VerifyTime = time.Now().Unix()
				n.Google = util.RequestGoogleBySocks5(i.ProxyIp)
				storage.StoreToSuccessDB(n)
			}
			storage.StoreToUnableDB(i)
			return
		}

		n := model.SuccessItem{ProxyIp: i.ProxyIp, Source: i.Source}
		n.VerifyTime = time.Now().Unix()
		n.Http = http
		n.Https = https
		n.Anonymous = anonymous
		if https {
			n.Google = util.RequestGoogleByHttps(i.ProxyIp)
		}
		// insert
		storage.StoreToSuccessDB(n)
		// remove
		storage.StoreToUnableDB(i)
	}()
}
